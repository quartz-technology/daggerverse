package docker

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/codebase/dockercompose"
	"dagger.io/dockersdk/module/object"
	"dagger.io/dockersdk/utils"
)

type serviceFunc struct {
	d       *Docker
	service *dockercompose.Service
}

//
// Name of the secret is lost when set
// func X(my-secret *dagger.Secret)

// Ask Marcos
// Set the entrypoint for when as-service will be called
func (s *serviceFunc) container(
	ctx context.Context,
	source *dockercompose.Source,
	env map[string]string,
	secretsEnv map[string]*dagger.Secret,
	mountedSecrets map[string]*dagger.Secret,
	mountedVolumes map[string]*dagger.Directory,
	caches map[string]string,
) (*dagger.Container, error) {
	var ctr *dagger.Container

	switch source.Type {
	case dockercompose.SourceTypeImage:
		ctr = dag.Container().From(source.Image.Ref)
	case dockercompose.SourceTypeDockerfile:
		buildOpts := dagger.DirectoryDockerBuildOpts{
			Dockerfile: source.Dockerfile.Dockerfile,
		}

		if source.Dockerfile.Target != nil {
			buildOpts.Target = *source.Dockerfile.Target
		}

		for key, value := range source.Dockerfile.BuildArgs {
			buildOpts.BuildArgs = append(buildOpts.BuildArgs, dagger.BuildArg{
				Name:  key,
				Value: *value,
			})
		}

		ctr = s.d.Dir.Directory(source.Dockerfile.Context).DockerBuild(buildOpts)
	default:
		return nil, fmt.Errorf("unknown source type %s", source.Type)
	}

	user, err := ctr.User(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if workdir := s.service.Workdir(); workdir != "" {
		ctr = ctr.WithWorkdir(workdir)
	}

	for key, value := range env {
		ctr = ctr.WithEnvVariable(key, value, dagger.ContainerWithEnvVariableOpts{Expand: true})
	}

	for name, secret := range secretsEnv {
		ctr = ctr.WithSecretVariable(name, secret)
	}

	for name, secret := range mountedSecrets {
		ctr = ctr.WithMountedSecret(fmt.Sprintf("/run/secrets/%s", name), secret)
	}

	for _, port := range s.service.Ports() {
		ctr = ctr.WithExposedPort(port)
	}

	for path, volume := range mountedVolumes {
		ctr = ctr.WithMountedDirectory(path, volume, dagger.ContainerWithMountedDirectoryOpts{
			Owner: fmt.Sprintf("%s:%s", user, user),
		})
	}

	for name, target := range caches {
		ctr = ctr.WithMountedCache(target, dag.CacheVolume(name), dagger.ContainerWithMountedCacheOpts{
			Owner: fmt.Sprintf("%s:%s", user, user),
		})
	}

	return ctr, nil
}

func (s *serviceFunc) Invoke(ctx context.Context, state object.State, input object.InputArgs) (object.Result, error) {
	docker, err := s.d.load(state)
	if err != nil {
		return nil, fmt.Errorf("failed to load object state: %w", err)
	}

	envMap, secretsMap := s.service.Environment()
	mountedSecretsName := s.service.MountedSecrets()
	mountedVolumePaths, cachesPaths := s.service.Volumes()

	// The image may be overwritten by the user
	source := s.service.Source()
	if source.Type == dockercompose.SourceTypeImage {
		source.Image.Ref = utils.LoadArgument[string]("image", input)
	}

	env := map[string]string{}
	for key := range envMap {
		env[key] = utils.LoadArgument[string](key, input)
	}

	secrets := map[string]*dagger.Secret{}
	for _, name := range secretsMap {
		if input[name] != nil {
			cliSecret := utils.LoadSecretFromID([]byte(input[name]))

			secretValue, err := cliSecret.Plaintext(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to add secret value: %w", err)
			}

			secrets[name] = dag.SetSecret(name, secretValue)
		}
	}

	mountedSecrets := map[string]*dagger.Secret{}
	for _, name := range mountedSecretsName {
		if input[name] != nil {
			cliSecret := utils.LoadSecretFromID([]byte(input[name]))

			secretValue, err := cliSecret.Plaintext(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to add secret value: %w", err)
			}

			mountedSecrets[name] = dag.SetSecret(name, secretValue)
		}
	}

	volumes := map[string]*dagger.Directory{}
	for _, volumePath := range mountedVolumePaths {
		if input[volumePath.Name()] != nil {
			volumes[volumePath.Target()] = utils.LoadDirectoryFromID([]byte(input["dir"]))
		}
	}

	caches := map[string]string{}
	for _, cache := range cachesPaths {
		caches[cache.Name()] = cache.Path()
	}

	return (*serviceFunc).container(
		&serviceFunc{d: docker, service: s.service}, ctx,
		source, env, secrets,
		mountedSecrets, volumes, caches,
	)
}

func (s *serviceFunc) AddTypeDefToObject(ctx context.Context, mod *dagger.Module, object *dagger.TypeDef) (*dagger.Module, *dagger.TypeDef) {
	typedef := dag.
		Function(s.service.Name(), dag.TypeDef().WithObject("Container")).
		WithDescription(fmt.Sprintf("Create a %s service container", s.service.Name()))

	source := s.service.Source()

	if source.Type == dockercompose.SourceTypeImage {
		typedef = typedef.WithArg("image", dag.TypeDef().WithKind(dagger.TypeDefKindStringKind), dagger.FunctionWithArgOpts{
			DefaultValue: utils.LoadDefaultValue(source.Image.Ref),
			Description:  "Image to use for the service",
		})
	}

	/////
	// Add environment variables
	env, secrets := s.service.Environment()
	for key, value := range env {
		opts := dagger.FunctionWithArgOpts{
			Description: fmt.Sprintf("Set environment variable %s", key),
		}

		if value != nil {
			opts.DefaultValue = utils.LoadDefaultValue(value)
		}

		typedef = typedef.WithArg(key, dag.TypeDef().WithKind(dagger.TypeDefKindStringKind), opts)
	}

	for _, name := range secrets {
		typedef = typedef.WithArg(name, dag.TypeDef().WithObject("Secret"), dagger.FunctionWithArgOpts{
			Description: fmt.Sprintf("Set secret environment variable %s", name),
		})
	}

	/////
	// Add mounted secrets
	mountedSecretsName := s.service.MountedSecrets()
	for _, name := range mountedSecretsName {
		typedef = typedef.WithArg(name, dag.TypeDef().WithObject("Secret"), dagger.FunctionWithArgOpts{
			Description: fmt.Sprintf("Secret %s to mount", name),
		})
	}

	/////
	// Add mounted volumes
	mountedVolumesPaths, _ := s.service.Volumes()
	for _, volumePath := range mountedVolumesPaths {
		typedef = typedef.WithArg(
			volumePath.Name(),
			dag.TypeDef().WithObject("Directory"),
			dagger.FunctionWithArgOpts{
				DefaultPath: volumePath.Origin(),
				Description: fmt.Sprintf("Mount directory at %s", volumePath.Target()),
			},
		)
	}

	return mod, object.WithFunction(typedef)
}
