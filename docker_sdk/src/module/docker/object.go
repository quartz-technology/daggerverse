package docker

import (
	"context"
	"encoding/json"
	"fmt"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/codebase/dockercompose"
	"dagger.io/dockersdk/codebase/dockerfile"
	"dagger.io/dockersdk/module/compose"
	"dagger.io/dockersdk/module/object"
	"dagger.io/dockersdk/utils"
)

type Docker struct {
	// Dir is the directory associated with the Docker object.
	Dir *dagger.Directory

	// name is the identifier for the Docker object.
	name string

	// dockerfile represents the Dockerfile associated with this Docker object.
	dockerfile *dockerfile.Dockerfile

	// dockercomposeFile represents the Docker Compose file linked with this Docker object.
	dockercomposeFile *dockercompose.DockerCompose

	// funcMap is a map of function names to their corresponding implementation.
	funcMap map[string]object.Function
}

// New creates a new Docker object with the specified name.
func New(name string) *Docker {
	return &Docker{
		name:       name,
		dockerfile: &dockerfile.Dockerfile{},
		funcMap:    map[string]object.Function{},
	}
}

// Name returns the name of the Docker object.
func (d *Docker) Name() string {
	return d.name
}

// Description provides a brief description of the Docker object.
func (d *Docker) Description() string {
	return "Execute docker function"
}

// AddTypeDef adds Docker function definitions to a module.
func (d *Docker) AddTypeDef(ctx context.Context) dagger.WithModuleFunc {
	return func(mod *dagger.Module) *dagger.Module {
		object := dag.TypeDef().WithObject(d.name)

		for _, fct := range d.funcMap {
			mod, object = fct.AddTypeDefToObject(ctx, mod, object)
		}

		return mod.WithObject(object)
	}
}

// New initializes a new Docker object using InputArgs.
func (d *Docker) New(input object.InputArgs) object.Object {
	var dir *dagger.Directory

	if input["dir"] != nil {
		dir = utils.LoadDirectoryFromID([]byte(input["dir"]))
	}

	return &Docker{
		Dir: dir,
	}
}

// Load reconstructs a Docker object from its State.
func (d *Docker) Load(state object.State) (object.Object, error) {
	return d.load(state)
}

// load helps in copying and returning a new Docker object from its state.
func (d *Docker) load(state object.State) (*Docker, error) {
	parentMap := make(map[string]interface{})
	err := json.Unmarshal(state, &parentMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal parent object: %w", err)
	}

	cpyDocker := &Docker{
		name:              d.name,
		dockerfile:        d.dockerfile,
		dockercomposeFile: d.dockercomposeFile,
		funcMap:           d.funcMap,
	}

	if parentMap["Dir"] != nil {
		cpyDocker.Dir = dag.LoadDirectoryFromID(dagger.DirectoryID(parentMap["Dir"].(string)))
	}

	return cpyDocker, nil
}

// Invoke calls a specific function from funcMap with the provided parameters.
func (d *Docker) Invoke(ctx context.Context, state object.State, fnName string, input object.InputArgs) (object.Result, error) {
	if d.funcMap[fnName] == nil {
		return nil, fmt.Errorf("unknown function %s", fnName)
	}

	return d.funcMap[fnName].Invoke(ctx, state, input)
}

// WithDockerfile associates a Dockerfile with the Docker object and adds
// a "Build" function.
func (d *Docker) WithDockerfile(dockerfile *dockerfile.Dockerfile) *Docker {
	d.dockerfile = dockerfile
	d.funcMap["Build"] = &buildFunc{d: d}

	return d
}

// WithDockerCompose associates a DockerCompose file with the Docker object and
// adds a "Compose" function.
func (d *Docker) WithDockerCompose(dockercomposeFile *dockercompose.DockerCompose) *Docker {
	d.dockercomposeFile = dockercomposeFile
	d.funcMap["Compose"] = &composeFunc{d: d}

	return d
}

// Deps returns a map of dependent objects for the Docker object.
func (d *Docker) Deps() map[string]object.Object {
	deps := make(map[string]object.Object)

	if d.dockercomposeFile != nil {
		composeObj := compose.New(d.Dir, d.dockercomposeFile)

		deps[composeObj.Name()] = composeObj
	}

	return deps
}
