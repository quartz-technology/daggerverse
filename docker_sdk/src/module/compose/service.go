package compose

import (
	"context"
	"fmt"

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

	// If as dep is true, the service will prefix the service name before
	// looking for the input arguments
	asDep bool
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
	mountedFiles map[string]*dagger.File,
	caches map[string]string,
	dependentServices []*proxy.Service,
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

	// Get exposed ports
	exposedPorts, err := ctr.ExposedPorts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get exposed ports: %w", err)
	}

	for _, port := range exposedPorts {
		exposedPort, err := port.Port(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get exposed port: %w", err)
		}

		ctr = ctr.WithExposedPort(exposedPort)

		s.service.WithExposedPort(exposedPort)
	}

	for path, volume := range mountedVolumes {
		ctr = ctr.WithMountedDirectory(path, volume, dagger.ContainerWithMountedDirectoryOpts{
			Owner: fmt.Sprintf("%s:%s", user, user),
		})
	}

	for path, file := range mountedFiles {
		ctr = ctr.WithMountedFile(path, file, dagger.ContainerWithMountedFileOpts{
			Owner: fmt.Sprintf("%s:%s", user, user),
		})
	}

	for name, target := range caches {
		ctr = ctr.WithMountedCache(target, dag.CacheVolume(name), dagger.ContainerWithMountedCacheOpts{
			Owner: fmt.Sprintf("%s:%s", user, user),
		})
	}

	entrypoint, exist := s.service.Entrypoint()
	if exist {
		ctr = ctr.WithEntrypoint(entrypoint)
	}

	command, exist := s.service.Command()
	if exist {
		ctr = ctr.WithDefaultArgs(command)
	}

	for _, service := range dependentServices {
		ctr = ctr.
			WithServiceBinding(service.Name, service.Service)
	}

	return ctr.Sync(ctx)
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

	service := &proxy.Service{
		Service: ctr.AsService(dagger.ContainerAsServiceOpts{UseEntrypoint: true}),
		Name:    s.service.Name(),
		Alias:   s.service.ContainerName(),
		Exposed: false,
	}

	ports := s.service.Ports()
	if len(ports) == 0 {
		return service, nil
	}

	service.Frontend = ports[0]
	service.Backend = ports[0]
	service.Exposed = true

	return service, nil
}

func (s *serviceFunc) formatInputArgName(argName string) string {
	if s.asDep {
		return fmt.Sprintf("%s_%s", s.service.Name(), argName)
	}

	return argName
}

func (s *serviceFunc) Invoke(ctx context.Context, state object.State, input object.InputArgs) (object.Result, error) {
	fmt.Printf("Invoking service %s; input: %#v\n", s.service.Name(), input)

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
		source.Image.Ref = utils.LoadArgument[string](s.formatInputArgName("image"), input)
	}

	env := map[string]string{}
	for key := range envMap {
		env[key] = utils.LoadArgument[string](s.formatInputArgName(utils.FormatEnvVariableName(key)), input)
	}

	secrets := map[string]*dagger.Secret{}
	for _, name := range secretsMap {
		if input[s.formatInputArgName(name)] != nil {
			cliSecret := utils.LoadSecretFromID([]byte(input[s.formatInputArgName(utils.FormatEnvVariableName(name))]))

			secretValue, err := cliSecret.Plaintext(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to add secret value: %w", err)
			}

			secrets[name] = dag.SetSecret(name, secretValue)
		} else {
			secrets[name] = dag.SetSecret(name, `""`)
		}
	}

	mountedSecrets := map[string]*dagger.Secret{}
	for _, name := range mountedSecretsName {
		if input[s.formatInputArgName(name)] != nil {
			cliSecret := utils.LoadSecretFromID([]byte(input[s.formatInputArgName(name)]))

			secretValue, err := cliSecret.Plaintext(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to add secret value: %w", err)
			}

			mountedSecrets[name] = dag.SetSecret(name, secretValue)
		}
	}

	volumes := map[string]*dagger.Directory{}
	mountedFiles := map[string]*dagger.File{}
	for _, volumePath := range mountedVolumePaths {
		if input[s.formatInputArgName(volumePath.Name())] != nil {
			if volumePath.IsDir() {
				volumes[volumePath.Target()] = utils.LoadDirectoryFromID([]byte(input[s.formatInputArgName(volumePath.Name())]))
			} else {
				mountedFiles[volumePath.Target()] = utils.LoadFileFromID([]byte(input[s.formatInputArgName(volumePath.Name())]))
			}
		}
	}

	caches := map[string]string{}
	for _, cache := range cachesPaths {
		caches[cache.Name()] = cache.Path()
	}

	dependentServices := []*proxy.Service{}
	for _, dependentServiceName := range s.service.DependsOn() {
		if compose.runningServices[dependentServiceName] != nil {
			dependentService := compose.runningServices[dependentServiceName]
			fmt.Printf("service %s is already running ; binding it to the dependent service %s\n", dependentServiceName, dependentService.Name, s.service.Name())

			dependentServices = append(dependentServices, dependentService)

			continue
		}

		fmt.Printf("service %s is not running yet; starting it\n", dependentServiceName)

		dockerComposeService, err := s.c.dockercompose.GetService(dependentServiceName)
		if err != nil {
			return nil, fmt.Errorf("failed to get service %s that %s depends on", dependentServiceName, s.service.Name())
		}

		serviceFct := &serviceFunc{c: compose, service: dockerComposeService, asDep: true}

		service, err := serviceFct.ToService(ctx, state, input)
		if err != nil {
			return nil, fmt.Errorf("failed to get dependent service: %w", err)
		}

		dependentServices = append(dependentServices, service)
	}

	fmt.Printf("Starting service %s\n", s.service.Name())

	return (*serviceFunc).container(
		&serviceFunc{c: compose, service: s.service, asDep: s.asDep}, ctx,
		source, env, secrets,
		mountedSecrets, volumes, mountedFiles,
		caches, dependentServices,
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

		if value != nil && *value != "" {
			opts.DefaultValue = utils.LoadDefaultValue(value)
		}

		args = append(args, &object.FunctionArg{
			Name: utils.FormatEnvVariableName(key),
			Type: dag.TypeDef().
				WithKind(dagger.TypeDefKindStringKind).
				// Environment variables are optional and will default to an empty
				// string if not set.
				WithOptional(true),
			Opts: opts,
		})
	}

	/////
	// Add environment variables secrets
	for _, name := range secrets {
		args = append(args, &object.FunctionArg{
			Name: utils.FormatEnvVariableName(name),
			Type: dag.TypeDef().WithObject("Secret").WithOptional(true),
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
		if volumePath.IsDir() {
			args = append(args, &object.FunctionArg{
				Name: volumePath.Name(),
				Type: dag.TypeDef().WithObject("Directory"),
				Opts: dagger.FunctionWithArgOpts{
					DefaultPath: volumePath.Origin(),
					Description: fmt.Sprintf("Mount directory at %s", volumePath.Target()),
				},
			})
		} else {
			args = append(args, &object.FunctionArg{
				Name: volumePath.Name(),
				Type: dag.TypeDef().WithObject("File"),
				Opts: dagger.FunctionWithArgOpts{
					DefaultPath: volumePath.Origin(),
					Description: fmt.Sprintf("Mount file at %s", volumePath.Target()),
				},
			})
		}
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
