package docker

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/utils"
)

func (d *Docker) build(platform *dagger.Platform, target *string, dockerfile *string, buildArgs []dagger.BuildArg, secrets []*dagger.Secret) *dagger.Container {
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

	return d.Dir.DockerBuild(opts)
}

func (d *Docker) buildFctTypeDef(ctx context.Context) (*dagger.Function, *dagger.TypeDef) {
	typedef := dag.Function("Build", dag.TypeDef().WithObject("Container")).
		WithDescription("Build a container from the Dockerfile in the current directory").
		WithArg("dockerfile",
			dag.TypeDef().WithKind(dagger.TypeDefKindStringKind).WithOptional(true),
			dagger.FunctionWithArgOpts{
				DefaultValue: utils.LoadDefaultValue(d.dockerfile.Filename()),
				Description:  "Path to the Dockerfile to use.",
			})

	/////
	// Add the build arguments
	for key, value := range d.dockerfile.Args() {
		buildArgOpts := dagger.FunctionWithArgOpts{
			Description: fmt.Sprintf("Set %s build argument", key),
		}

		if value != "" {
			buildArgOpts.DefaultValue = utils.LoadDefaultValue(value)
		}

		typedef = typedef.WithArg(key, dag.TypeDef().WithKind(dagger.TypeDefKindStringKind), buildArgOpts)
	}

	//////
	// Add the secrets arguments
	for _, secret := range d.dockerfile.Secrets() {
		typedef = typedef.WithArg(secret,
			dag.TypeDef().WithObject("Secret"),
			dagger.FunctionWithArgOpts{
				Description: fmt.Sprintf("Set %s secret", secret),
			})
	}

	/////
	// Add the platform argument
	defaultPlatformArgOpts := dagger.FunctionWithArgOpts{
		Description: "Platform to build.",
	}

	// Ideally, we want to get the default platform from the host but we don't want to
	// fail the call if we can't get it.
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

	/////
	// If stages are declared in the Dockerfile, we add an enum and the stage option to the function.
	if len(d.dockerfile.Stages()) != 0 {
		stageTypeDef := dag.TypeDef().WithEnum(fmt.Sprintf("%sStage", d.name))

		for _, stage := range d.dockerfile.Stages() {
			stageTypeDef = stageTypeDef.WithEnumValue(stage)
		}

		typedef = typedef.
			WithArg("target", dag.TypeDef().WithEnum(fmt.Sprintf("%sStage", d.name)).
				WithOptional(true),
				dagger.FunctionWithArgOpts{
					Description: "Target stage to build.",
				},
			)

		return typedef, stageTypeDef
	}

	return typedef, nil
}
