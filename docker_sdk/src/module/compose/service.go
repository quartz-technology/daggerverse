package compose

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/codebase/dockercompose"
	"dagger.io/dockersdk/module/object"
	"dagger.io/dockersdk/module/proxy"
	"dagger.io/dockersdk/utils"
)

type serviceFunc struct {
	c       *Compose
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
	dependentServices map[string]*dagger.Container,
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

		ctr = s.c.Dir.Directory(source.Dockerfile.Context).DockerBuild(buildOpts)
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

	for name, service := range dependentServices {
		ctr = ctr.WithServiceBinding(name, service.AsService())
	}

	return ctr, nil
}

func (s *serviceFunc) ToContainer(ctx context.Context, state object.State, input object.InputArgs) (*dagger.Container, error) {
	ctrRes, err := s.Invoke(ctx, state, input)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke service %s: %w", s.service.Name(), err)
	}

	ctr, ok := ctrRes.(*dagger.Container)
	if !ok {
		return nil, fmt.Errorf("expected container result, got %T", ctrRes)
	}

	return ctr, nil
}

func (s *serviceFunc) ToService(ctx context.Context, state object.State, input object.InputArgs) (*proxy.Service, error) {
	ctr, err := s.ToContainer(ctx, state, input)
	if err != nil {
		return nil, fmt.Errorf("failed to convert service %s to container: %w", s.service.Name(), err)
	}

	ports := s.service.Ports()
	if len(ports) == 0 {
		return nil, fmt.Errorf("service %s has no ports", s.service.Name())
	}

	return &proxy.Service{
		Service:  ctr.AsService(dagger.ContainerAsServiceOpts{UseEntrypoint: true}),
		Name:     s.service.Name(),
		Frontend: ports[0],
		Backend:  ports[0],
	}, nil
}

func (s *serviceFunc) Invoke(ctx context.Context, state object.State, input object.InputArgs) (object.Result, error) {
	compose, err := s.c.load(state)
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
			volumes[volumePath.Target()] = utils.LoadDirectoryFromID([]byte(input[volumePath.Name()]))
		}
	}

	caches := map[string]string{}
	for _, cache := range cachesPaths {
		caches[cache.Name()] = cache.Path()
	}

	dependentServices := map[string]*dagger.Container{}
	for _, dependentServiceName := range s.service.DependsOn() {
		dockerComposeService, err := s.c.dockercompose.GetService(dependentServiceName)
		if err != nil {
			return nil, fmt.Errorf("failed to get service %s that %s depends on", dependentServiceName, s.service.Name())
		}

		service := &serviceFunc{c: compose, service: dockerComposeService}
		servicePrefix := fmt.Sprintf("%s_", dependentServiceName)

		serviceInput := input
		for argName, argValue := range input {
			if strings.HasPrefix(argName, servicePrefix) {
				serviceInput[strings.TrimPrefix(argName, servicePrefix)] = argValue
			}
		}

		ctr, err := service.ToContainer(ctx, state, serviceInput)
		if err != nil {
			return nil, fmt.Errorf("failed to get dependent service: %w", err)
		}

		dependentServices[dependentServiceName] = ctr
	}

	return (*serviceFunc).container(
		&serviceFunc{c: compose, service: s.service}, ctx,
		source, env, secrets,
		mountedSecrets, volumes, caches,
		dependentServices,
	)
}

func (s *serviceFunc) Arguments() []*object.FunctionArg {
	args := []*object.FunctionArg{}

	/////
	// Add image if necessary
	source := s.service.Source()
	if source.Type == dockercompose.SourceTypeImage {
		args = append(args, &object.FunctionArg{
			Name: "image",
			Type: dag.TypeDef().WithKind(dagger.TypeDefKindStringKind),
			Opts: dagger.FunctionWithArgOpts{
				DefaultValue: utils.LoadDefaultValue(source.Image.Ref),
				Description:  "Image to use for the service",
			},
		})
	}

	//////
	// Add environment variables
	env, secrets := s.service.Environment()
	for key, value := range env {
		opts := dagger.FunctionWithArgOpts{
			Description: fmt.Sprintf("Set environment variable %s", key),
		}

		if value != nil {
			opts.DefaultValue = utils.LoadDefaultValue(value)
		}

		args = append(args, &object.FunctionArg{
			Name: key,
			Type: dag.TypeDef().WithKind(dagger.TypeDefKindStringKind),
			Opts: opts,
		})
	}

	/////
	// Add environment variables secrets
	for _, name := range secrets {
		args = append(args, &object.FunctionArg{
			Name: name,
			Type: dag.TypeDef().WithObject("Secret"),
			Opts: dagger.FunctionWithArgOpts{
				Description: fmt.Sprintf("Set secret environment variable %s", name),
			},
		})
	}

	mountedSecretsName := s.service.MountedSecrets()
	for _, name := range mountedSecretsName {
		args = append(args, &object.FunctionArg{
			Name: name,
			Type: dag.TypeDef().WithObject("Secret"),
			Opts: dagger.FunctionWithArgOpts{
				Description: fmt.Sprintf("Secret %s to mount", name),
			},
		})
	}

	/////
	// Add mounted volumes
	mountedVolumesPaths, _ := s.service.Volumes()
	for _, volumePath := range mountedVolumesPaths {
		args = append(args, &object.FunctionArg{
			Name: volumePath.Name(),
			Type: dag.TypeDef().WithObject("Directory"),
			Opts: dagger.FunctionWithArgOpts{
				DefaultPath: volumePath.Origin(),
				Description: fmt.Sprintf("Mount directory at %s", volumePath.Target()),
			},
		})
	}

	return args
}

func (s *serviceFunc) AddTypeDefToObject(ctx context.Context, mod *dagger.Module, object *dagger.TypeDef) (*dagger.Module, *dagger.TypeDef) {
	typedef := dag.
		Function(s.service.Name(), dag.TypeDef().WithObject("Container")).
		WithDescription(fmt.Sprintf("Create a %s service container", s.service.Name()))

	args := s.Arguments()
	for _, arg := range args {
		typedef = typedef.WithArg(arg.Name, arg.Type, arg.Opts)
	}

	/////
	// Add depends on service argument

	for _, dependencyName := range s.service.DependsOn() {
		service, exist := s.c.funcMap[dependencyName]
		if !exist {
			panic(fmt.Errorf("service %s does not exist but %s depends on it", dependencyName, s.service.Name()))
		}

		serviceArgs := service.Arguments()
		for _, arg := range serviceArgs {
			typedef = typedef.WithArg(
				fmt.Sprintf("%s_%s", dependencyName, arg.Name),
				arg.Type,
				arg.Opts,
			)
		}
	}

	return mod, object.WithFunction(typedef)
}
