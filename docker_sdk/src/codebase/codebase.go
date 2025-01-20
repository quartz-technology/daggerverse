package codebase

import (
	"fmt"

	"dagger.io/dockersdk/dockerfile"
)

type Codebase struct {
	dockerfile *dockerfile.Dockerfile
}

func New() (*Codebase, error) {
	dockerfile, exists, err := getDockerfile(CodebasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get Dockerfile: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("Dockerfile not found in %s", CodebasePath)
	}

	// TODO: Check for sub directories/dockerfile

	return &Codebase{dockerfile: dockerfile}, nil
}

func (c *Codebase) Dockerfile() *dockerfile.Dockerfile {
	return c.dockerfile
}