package codebase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dockersdk/codebase/dockercompose"
	"dagger.io/dockersdk/codebase/dockerfile"
	"dagger.io/dockersdk/codebase/finder"
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

	finder := finder.New(CodebasePath, dir)

	dockerfile, dockerfileExists, err := getDockerfile(finder)
	if err != nil {
		return nil, fmt.Errorf("failed to get Dockerfile: %w", err)
	}

	dockercompose, composeExistsExists, err := getDockerCompose(ctx, finder)
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

func getDockerfile(finder *finder.Finder) (*dockerfile.Dockerfile, bool, error) {
	patterns := []string{"Dockerfile", "*.Dockerfile"}

	dockerfilePath, exist := finder.FindFileFromPattern(patterns)
	if !exist {
		return nil, false, nil
	}

	file, err := os.Open(dockerfilePath)
	if err != nil {
		return nil, true, fmt.Errorf("failed to open Dockerfile: %w", err)
	}
	defer file.Close()

	filename := filepath.Base(dockerfilePath)
	dockerfile, err := dockerfile.NewDockerfile(filename, file)
	if err != nil {
		return nil, true, fmt.Errorf("failed to parse Dockerfile: %w", err)
	}

	return dockerfile, true, nil
}

func getDockerCompose(ctx context.Context, finder *finder.Finder) (*dockercompose.DockerCompose, bool, error) {
	patterns := []string{"docker-compose.yml", "docker-compose.yaml", "compose.yaml", "compose.yml"}

	dockerComposePath, exist := finder.FindFileFromPattern(patterns)
	if !exist {
		return nil, false, nil
	}

	filename := filepath.Base(dockerComposePath)
	fileContent, err := os.ReadFile(dockerComposePath)
	if err != nil {
		return nil, true, fmt.Errorf("failed to get %s content: %w", filename, err)
	}

	compose, err := dockercompose.NewDockerCompose(ctx, filename, fileContent, finder)
	if err != nil {
		return nil, true, fmt.Errorf("failed to parse docker-compose.yml: %w", err)
	}

	return compose, true, nil
}