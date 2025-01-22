package codebase

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dockersdk/codebase/dockercompose"
	"dagger.io/dockersdk/codebase/dockerfile"
)

const CodebasePath = "/app"


type Codebase struct {
	dockerfile    *dockerfile.Dockerfile
	dockercompose *dockercompose.DockerCompose
}

func New(ctx context.Context) (*Codebase, error) {
	dir, err := os.ReadDir(CodebasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read source directory: %w", err)
	}

	dockerfile, dockerfileExists, err := getDockerfile(CodebasePath, dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get Dockerfile: %w", err)
	}

	dockercompose, composeExistsExists, err := getDockerCompose(ctx, CodebasePath, dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get docker-compose file: %w", err)
	}

	if !dockerfileExists && !composeExistsExists {
		return nil, fmt.Errorf("Dockerfile or docker-compose.yml not found in user project")
	}

	// TODO: Check for sub directories/dockerfile

	return &Codebase{
		dockerfile:    dockerfile,
		dockercompose: dockercompose,
	}, nil
}