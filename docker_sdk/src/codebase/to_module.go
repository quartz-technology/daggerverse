package codebase

import (
	"dagger.io/dockersdk/module"
	"dagger.io/dockersdk/module/docker"
)

// Converts a Codebase instance to a Dagger Docker module.
//
// Initializes a new Docker module with the given name and optionally configures
// it with a Dockerfile and a Docker Compose file if they are present in the
// Codebase instance.
func (c *Codebase) ToModule(name string) *module.Module {
	dockerModule := docker.New("Docker")

	if c.dockerfile != nil {
		dockerModule = dockerModule.WithDockerfile(c.dockerfile)
	}

	if c.dockercompose != nil {
		dockerModule = dockerModule.WithDockerCompose(c.dockercompose)
	}

	return module.Build(name, dockerModule)
}
