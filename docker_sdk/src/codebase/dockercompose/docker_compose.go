package dockercompose

import (
	"context"
	"fmt"

	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
)

type DockerCompose struct {
	filename string
	project  *types.Project
}

func NewDockerCompose(ctx context.Context, filename string, content []byte) (*DockerCompose, error) {
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
	}, nil
}

func (d *DockerCompose) Services() []*Service {
  services := make([]*Service, len(d.project.Services))
	
	for i, service := range d.project.Services {
		services[i] = NewService(&service)
	}

	return services
}

func (d *DockerCompose) String() string {
	yaml, err := d.project.MarshalYAML()
	if err != nil {
		return "could not display docker-compose.yml"
	}

	return string(yaml)
}
