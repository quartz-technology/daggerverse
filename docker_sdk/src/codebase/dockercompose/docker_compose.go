package dockercompose

import (
	"context"
	"fmt"

	"dagger.io/dockersdk/codebase/finder"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
)

type DockerCompose struct {
	filename string
	project  *types.Project

	finder *finder.Finder
}

func NewDockerCompose(ctx context.Context, filename string, content []byte, finder *finder.Finder) (*DockerCompose, error) {
	project, err := loader.LoadWithContext(ctx, types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{
				Config: map[string]interface{}{
					"name": "dockersdk",
				},
			},
			{
				Content: content,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load %s: %w", filename, err)
	}

	return &DockerCompose{
		filename: filename,
		project:  project,
		finder: finder,
	}, nil
}

func (d *DockerCompose) Services() []*Service {
  services := make([]*Service, len(d.project.Services))
	
	for i, service := range d.project.Services {
		services[i] = NewService(d, &service, d.finder)
	}

	return services
}

func (d *DockerCompose) GetService(name string) (*Service, error) {
	for _, service := range d.Services() {
		if service.Name() == name {
			return service, nil
		}
	}

	return nil, fmt.Errorf("no such service: %s", name)
}

func (d *DockerCompose) String() string {
	yaml, err := d.project.MarshalYAML()
	if err != nil {
		return "could not display docker-compose.yml"
	}

	return string(yaml)
}
