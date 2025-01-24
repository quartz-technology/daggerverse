package codebase

import (
	"dagger.io/dockersdk/module"
	"dagger.io/dockersdk/module/docker"
)

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
