package docker

import (
	"context"
	"encoding/json"
	"fmt"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/codebase"
	"dagger.io/dockersdk/dockerfile"
	"dagger.io/dockersdk/integrations/object"
	"dagger.io/dockersdk/utils"
)

type Docker struct {
	Dir *dagger.Directory

	name string

	// Specific information retrieved from the Dockerfile
	dockerfile *dockerfile.Dockerfile
}

func New(name string, code *codebase.Codebase) object.Object {
	return &Docker{
		name:       name,
		dockerfile: code.Dockerfile(),
	}
}

func (d *Docker) Name() string {
	return d.name
}

func (d *Docker) Description() string {
	return "Execute docker function"
}

func (d *Docker) AddTypeDef(ctx context.Context) dagger.WithModuleFunc {
	return func(mod *dagger.Module) *dagger.Module {
		fct, enum := d.buildFctTypeDef(ctx)
		if enum != nil {
			mod = mod.WithEnum(enum)
		}

		object := dag.
			TypeDef().
			WithObject(d.name).
			WithFunction(fct)

		return mod.WithObject(object)
	}
}

func (d *Docker) Load(state object.State) (object.Object, error) {
	parentMap := make(map[string]interface{})
	err := json.Unmarshal(state, &parentMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal parent object: %w", err)
	}

	cpyDocker := &Docker{name: d.name, dockerfile: d.dockerfile}

	if parentMap["Dir"] != nil {
		cpyDocker.Dir = dag.LoadDirectoryFromID(dagger.DirectoryID(parentMap["Dir"].(string)))
	}

	return cpyDocker, nil
}

func (d *Docker) New(input object.InputArgs) object.Object {
	var dir *dagger.Directory

	if input["dir"] != nil {
		dir = utils.LoadDirectoryFromID([]byte(input["dir"]))
	}

	return &Docker{
		Dir: dir,
	}
}

func (d *Docker) Invoke(ctx context.Context, fnName string, input object.InputArgs) (object.Result, error) {
	switch fnName {
	case "Build":
		platform := utils.LoadArgument[dagger.Platform]("platform", input)
		target := utils.LoadArgument[string]("target", input)
		dockerfile := utils.LoadArgument[string]("dockerfile", input)

		buildArgs := []dagger.BuildArg{}
		for key := range d.dockerfile.Args() {
			if input[key] != nil {
				buildArgs = append(buildArgs, dagger.BuildArg{
					Name:  key,
					Value: utils.LoadArgument[string](key, input),
				})
			}
		}

		// To load secret we need to load the value first and then reassign a secret
		// with the right value since it's obfuscated by the CLI.
		// TODO: find a better way to do this.
		secrets := []*dagger.Secret{}
		for _, secretKey := range d.dockerfile.Secrets() {
			if input[secretKey] != nil {
				cliSecret := utils.LoadSecretFromID([]byte(input[secretKey]))

				secretValue, err := cliSecret.Plaintext(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to add secret value: %w", err)
				}

				secrets = append(secrets, dag.SetSecret(secretKey, secretValue))
			}
		}

		return d.build(&platform, &target, &dockerfile, buildArgs, secrets), nil
	default:
		return nil, fmt.Errorf("unknown function %s", fnName)
	}
}
