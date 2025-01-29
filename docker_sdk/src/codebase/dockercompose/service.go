package dockercompose

import (
	"fmt"
	"strconv"

	"dagger.io/dockersdk/codebase/finder"
	"dagger.io/dockersdk/utils"
	"github.com/compose-spec/compose-go/types"
)

type SourceImage struct {
	Ref string
}

type SourceDockerfile struct {
	Context    string
	Dockerfile string
	BuildArgs  map[string]*string
	Target     *string
	// TODO(TomChv): We do not support secrets in inline compose build to avoid
	// DX conflict if secrets are shared between run & build.
	// Secrets    []string
}

type SourceType string

const (
	SourceTypeImage      SourceType = "image"
	SourceTypeDockerfile SourceType = "dockerfile"
)

type Source struct {
	Type       SourceType
	Image      *SourceImage
	Dockerfile *SourceDockerfile
}

type Service struct {
	sourceCompose *DockerCompose
	s             *types.ServiceConfig
	finder *finder.Finder

	exposedPorts []int
}

func NewService(sourceCompose *DockerCompose, service *types.ServiceConfig, finder *finder.Finder) *Service {
	return &Service{
		sourceCompose: sourceCompose,
		s: service,
		finder: finder,
	}
}

func (s *Service) Name() string {
	return s.s.Name
}

func (s *Service) ContainerName() string {
	if s.s.ContainerName != "" {
		return s.s.ContainerName
	}

	return s.s.Name
}

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

func (s *Service) Image() string {
	return s.s.Image
}

func (s *Service) Dockerfile() string {
	return s.s.Build.Dockerfile
}

func (s *Service) Workdir() string {
	return s.s.WorkingDir
}

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

	// Clean duplicates
	ports = utils.RemoveListDuplicates(ports)

	return ports
}

func (s *Service) WithExposedPort(port int) *Service {
	s.exposedPorts = append(s.exposedPorts, port)

	return s
}

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

func (s *Service) MountedSecrets() []string {
	secrets := []string{}

	for _, secret := range s.s.Secrets {
		secrets = append(secrets, secret.Source)
	}

	return secrets
}

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

			// If we can't get the volumes, we'll convert it to cache
			if err != nil {
				caches = append(caches, &Cache{name: source, path: v.Target})

				continue
			}

			volumes = append(volumes, &Volume{origin: source, target: v.Target, isDir: isDir})
		}
	}

	return volumes, caches
}

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

func (s *Service) Entrypoint() ([]string, bool) {
	if s.s.Entrypoint != nil {
		return s.s.Entrypoint, true
	}

	return nil, false
}

func (s *Service) Command() ([]string, bool) {
	if s.s.Command != nil {
		return s.s.Command, true
	}

	return nil, false
}