package codebase

import (
	"fmt"

	"dagger.io/dockersdk/module"
	"dagger.io/dockersdk/module/docker"
)

func (c *Codebase) ToModule(name string) *module.Module {
	dockerModule := docker.New("Docker")

	if c.dockerfile != nil {
		dockerModule = dockerModule.WithDockerfile(c.dockerfile)
	}

	if c.dockercompose != nil {

		for _, service := range c.dockercompose.Services() {
			dockerModule = dockerModule.WithService(service)
		}
	}

	return module.Build(name, dockerModule)
}
