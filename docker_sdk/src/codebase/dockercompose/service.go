package dockercompose

import (
	"fmt"
	"strconv"

	"dagger.io/dockersdk/codebase/finder"
	"dagger.io/dockersdk/utils"
	"github.com/compose-spec/compose-go/types"
)

// Service represents a Docker Compose service configuration.
type Service struct {
	// sourceCompose is the DockerCompose object associated with the service.
	sourceCompose *DockerCompose

	// s is a reference to the parsed service configuration.
	s *types.ServiceConfig

	// finder assists in finding files in the host's directory.
	finder *finder.Finder

	// exposedPorts holds additional ports exposed by the service's image.
	//
	// These can only be retrieved after the service's image has been pulled.
	exposedPorts []int
}

func NewService(sourceCompose *DockerCompose, service *types.ServiceConfig, finder *finder.Finder) *Service {
	return &Service{
		sourceCompose: sourceCompose,
		s:             service,
		finder:        finder,
	}
}

// Name returns the name of the service.
func (s *Service) Name() string {
	return s.s.Name
}

// ContainerName returns the container name if specified, else the service name.
func (s *Service) ContainerName() string {
	if s.s.ContainerName != "" {
		return s.s.ContainerName
	}
	return s.s.Name
}

// Source returns the source of the service, either image or Dockerfile.
//
// If the service is not defined in the Compose file, this will leads to a panic.
func (s *Service) Source() *Source {
	if s.s.Image != "" {
		return &Source{
			Type: SourceTypeImage,
			Image: &SourceImage{
				Ref: s.s.Image,
			},
		}
	}

	if s.s.Build != nil {
		dockerfile := &SourceDockerfile{
			Dockerfile: s.s.Build.Dockerfile,
			Context:    trimHostPath(s.s.Build.Context),
		}

		if s.s.Build.Args != nil {
			for key, value := range s.s.Build.Args {
				dockerfile.BuildArgs[key] = value
			}
		}

		if s.s.Build.Target != "" {
			dockerfile.Target = &s.s.Build.Target
		}

		return &Source{
			Type:       SourceTypeDockerfile,
			Dockerfile: dockerfile,
		}
	}

	panic(fmt.Sprintf("no source found in the service %s", s.s.Name))
}

// Workdir returns the working directory for the service.
func (s *Service) Workdir() string {
	return s.s.WorkingDir
}

// Ports returns all published and exposed ports for the service.
func (s *Service) Ports() []int {
	ports := []int{}

	for _, port := range s.s.Ports {
		published, err := strconv.Atoi(port.Published)
		if err != nil {
			fmt.Println(fmt.Errorf("failed to parse port published: %w, ignoring it", err))
		}

		ports = append(ports, published)
	}

	for _, port := range s.s.Expose {
		published, err := strconv.Atoi(port)
		if err != nil {
			fmt.Println(fmt.Errorf("failed to parse port published: %w, ignoring it", err))
		}

		ports = append(ports, published)
	}

	for _, port := range s.exposedPorts {
		ports = append(ports, port)
	}

	// Remove duplicates if there are any.
	ports = utils.RemoveListDuplicates(ports)

	return ports
}

// WithExposedPort adds an additional exposed port to the service.
func (s *Service) WithExposedPort(port int) *Service {
	s.exposedPorts = append(s.exposedPorts, port)
	return s
}

// Environment returns environment variables and secrets for the service.
func (s *Service) Environment() (env map[string]*string, secrets []string) {
	env = map[string]*string{}

	for key, value := range s.s.Environment {
		if value == nil {
			secrets = append(secrets, key)

			continue
		}

		env[key] = value
	}

	return env, secrets
}

// MountedSecrets lists secrets mounted in the service configuration.
func (s *Service) MountedSecrets() []string {
	secrets := []string{}

	for _, secret := range s.s.Secrets {
		secrets = append(secrets, secret.Source)
	}

	return secrets
}

// Volumes returns all the volumes and caches used by the service.
//
// If a volume is not defined in the host's directory, it will be transformed
// into a cache volume.
func (s *Service) Volumes() ([]*Volume, []*Cache) {
	volumes := []*Volume{}
	caches := []*Cache{}

	for _, v := range s.s.Volumes {
		switch v.Type {
		case "volume":
			caches = append(caches, &Cache{name: v.Source, path: v.Target})
		case "bind":
			source := trimHostPath(v.Source)

			isDir, err := s.finder.IsPathDirectory(source)
			if err != nil {
				caches = append(caches, &Cache{name: source, path: v.Target})

				continue
			}

			volumes = append(volumes, &Volume{origin: source, target: v.Target, isDir: isDir})
		}
	}

	return volumes, caches
}

// DependsOn retrieves a list of services this service depends on.
func (s *Service) DependsOn() []string {
	dependentServices := map[string]bool{}

	for key := range s.s.DependsOn {
		dependentServices[key] = true

		service, err := s.sourceCompose.GetService(key)
		if err != nil {
			fmt.Printf("failed to get service %s: %s\n", key, err.Error())

			continue
		}

		serviceDeps := service.DependsOn()
		for _, dep := range serviceDeps {
			dependentServices[dep] = true
		}
	}

	var dependentServicesList []string
	for service := range dependentServices {
		dependentServicesList = append(dependentServicesList, service)
	}

	return dependentServicesList
}

// Entrypoint returns the entrypoint of the service, if defined.
func (s *Service) Entrypoint() ([]string, bool) {
	if s.s.Entrypoint != nil {
		return s.s.Entrypoint, true
	}

	return nil, false
}

// Command returns the command executed in the container, if defined.
func (s *Service) Command() ([]string, bool) {
	if s.s.Command != nil {
		return s.s.Command, true
	}

	return nil, false
}
