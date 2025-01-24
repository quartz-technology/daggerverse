package codebase

import (
	"context"
	"fmt"
	"os"
	"strings"

	"dagger.io/dockersdk/codebase/dockercompose"
	"dagger.io/dockersdk/codebase/dockerfile"
)

func getDockerfile(dirPath string, dir []os.DirEntry) (*dockerfile.Dockerfile, bool, error) {
	for _, entry := range dir {
		if entry.Name() == "Dockerfile" || strings.HasSuffix(entry.Name(), ".Dockerfile") {
			file, err := os.Open(dirPath + "/" + entry.Name())
			if err != nil {
				return nil, true, fmt.Errorf("failed to open Dockerfile: %w", err)
			}
			defer file.Close()

			dockerfile, err := dockerfile.NewDockerfile(entry.Name(), file)
			if err != nil {
				return nil, true, fmt.Errorf("failed to parse Dockerfile: %w", err)
			}

			return dockerfile, true, nil
		}
	}

	return nil, false, nil
}

func getDockerCompose(ctx context.Context, dirPath string, dir []os.DirEntry) (*dockercompose.DockerCompose, bool, error) {
	isComposeFile := func(name string) bool {
		for _, composeFile := range []string{"docker-compose.yml", "docker-compose.yaml", "compose.yaml", "compose.yml"} {
			if name == composeFile {
				return true
			}
		}

		return false
	}

	for _, entry := range dir {
		if isComposeFile(entry.Name()) {
			fileContent, err := os.ReadFile(dirPath + "/" + entry.Name())
			if err != nil {
				return nil, true, fmt.Errorf("failed to get %s content: %w", entry.Name(), err)
			}

			compose, err := dockercompose.NewDockerCompose(ctx, entry.Name(), fileContent)
			if err != nil {
				return nil, true, fmt.Errorf("failed to parse docker-compose.yml: %w", err)
			}

			return compose, true, nil
		}
	}

	return nil, false, nil
}
