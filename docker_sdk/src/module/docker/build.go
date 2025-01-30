package docker

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/module/object"
	"dagger.io/dockersdk/utils"
)

// buildFunc encapsulates Docker build methods.
type buildFunc struct {
	// d holds the Docker object.
	d *Docker
}

// build constructs a container from the given parameters.
//
// It takes optional platform, target, Dockerfile path, build arguments, and
// secrets.
//
// Returns a built dagger.Container.
func (b *buildFunc) build(
	platform *dagger.Platform,
	target *string,
	dockerfile *string,
	buildArgs []dagger.BuildArg,
	secrets []*dagger.Secret,
) *dagger.Container {
	opts := dagger.DirectoryDockerBuildOpts{
		BuildArgs: buildArgs,
		Secrets:   secrets,
	}

	if platform != nil {
		opts.Platform = *platform
	}

	if target != nil {
		opts.Target = *target
	}

	if dockerfile != nil {
		opts.Dockerfile = *dockerfile
	}

	return b.d.Dir.DockerBuild(opts)
}

// Invoke executes the build process.
//
// It verifies the presence of a Dockerfile and loads necessary resources from the
// object state.
//
// It accepts context, object state, and input arguments, then performs
// the build using the specified options.
//
// Returns the build result or an error.
func (b *buildFunc) Invoke(ctx context.Context, state object.State, input object.InputArgs) (object.Result, error) {
	if b.d.dockerfile == nil {
		return nil, fmt.Errorf("Build function invoked before Dockerfile is set")
	}

	// Loads Docker object from object state.
	docker, err := b.d.load(state)
	if err != nil {
		return nil, fmt.Errorf("failed to load object state: %w", err)
	}

	// Loads platform, dockerfile and target from input.
	platform := utils.LoadArgument[dagger.Platform]("platform", input)
	target := utils.LoadArgument[string]("target", input)
	dockerfile := utils.LoadArgument[string]("dockerfile", input)

	// Loads build arguments from input.
	buildArgs := []dagger.BuildArg{}
	for key := range b.d.dockerfile.Args() {
		if input[key] != nil {
			buildArgs = append(buildArgs, dagger.BuildArg{
				Name:  key,
				Value: utils.LoadArgument[string](key, input),
			})
		}
	}

	// Adds secrets after decrypting with CLI utilities.
	//
	// This workaround is required since the secret's name
	// isn't the same as the identifier defined in the Dockerfile.
	secrets := []*dagger.Secret{}
	for _, secretKey := range b.d.dockerfile.Secrets() {
		if input[secretKey] != nil {
			cliSecret := utils.LoadSecretFromID([]byte(input[secretKey]))

			secretValue, err := cliSecret.Plaintext(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to add secret value: %w", err)
			}

			secrets = append(secrets, dag.SetSecret(secretKey, secretValue))
		}
	}

	return (*buildFunc).build(&buildFunc{d: docker}, &platform, &target, &dockerfile, buildArgs, secrets), nil
}

// Arguments is a placeholder method not invoked for this function
// required to implements the object.Function interface.
//
// This function should never be called for this function.
func (b *buildFunc) Arguments() []*object.FunctionArg {
	return nil
}

// AddTypeDefToObject adds the "Build" function definition to a Dagger
// module's object.
//
// It defines the function signature including Dockerfile path, build arguments,
// secrets, platform, and target stages.
//
// It returns the updated module and type definition.
func (b *buildFunc) AddTypeDefToObject(ctx context.Context, mod *dagger.Module, object *dagger.TypeDef) (*dagger.Module, *dagger.TypeDef) {
	typedef := dag.Function("Build", dag.TypeDef().WithObject("Container")).
		WithDescription("Build a container from the Dockerfile in the current directory").
		WithArg("dockerfile",
			dag.TypeDef().WithKind(dagger.TypeDefKindStringKind).WithOptional(true),
			dagger.FunctionWithArgOpts{
				DefaultValue: utils.LoadDefaultValue(b.d.dockerfile.Filename()),
				Description:  "Path to the Dockerfile to use.",
			})

	// Add the build arguments
	for key, value := range b.d.dockerfile.Args() {
		buildArgOpts := dagger.FunctionWithArgOpts{
			Description: fmt.Sprintf("Set %s build argument", key),
		}

		if value != "" {
			buildArgOpts.DefaultValue = utils.LoadDefaultValue(value)
		}

		typedef = typedef.WithArg(key, dag.TypeDef().WithKind(dagger.TypeDefKindStringKind), buildArgOpts)
	}

	// Add the secrets arguments
	for _, secret := range b.d.dockerfile.Secrets() {
		typedef = typedef.WithArg(secret,
			dag.TypeDef().WithObject("Secret"),
			dagger.FunctionWithArgOpts{
				Description: fmt.Sprintf("Set %s secret", secret),
			})
	}

	// Add the platform argument
	defaultPlatformArgOpts := dagger.FunctionWithArgOpts{
		Description: "Platform to build.",
	}

	defaultPlatform, err := dag.DefaultPlatform(ctx)
	if err == nil {
		defaultPlatformArgOpts.DefaultValue = utils.LoadDefaultValue(defaultPlatform)
	}

	typedef = typedef.
		WithArg("platform", dag.
			TypeDef().
			WithScalar("Platform").
			WithOptional(true),
			defaultPlatformArgOpts,
		)

	// Add target stage option if stages are declared in the Dockerfile.
	if len(b.d.dockerfile.Stages()) != 0 {
		stageTypeDef := dag.TypeDef().WithEnum(fmt.Sprintf("%sStage", b.d.name))

		for _, stage := range b.d.dockerfile.Stages() {
			stageTypeDef = stageTypeDef.WithEnumValue(stage)
		}

		typedef = typedef.
			WithArg("target", dag.TypeDef().WithEnum(fmt.Sprintf("%sStage", b.d.name)).
				WithOptional(true),
				dagger.FunctionWithArgOpts{
					Description: "Target stage to build.",
				},
			)

		mod = mod.WithEnum(stageTypeDef)
	}

	return mod, object.WithFunction(typedef)
}
