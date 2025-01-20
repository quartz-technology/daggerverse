package codebase

import (
	"fmt"
	"os"
	"strings"

	"dagger.io/dockersdk/dockerfile"
)

func getDockerfile(dirPath string) (*dockerfile.Dockerfile, bool, error) {
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range dir {
		if entry.Name() == "Dockerfile" || strings.HasSuffix(entry.Name(), ".Dockerfile") {
			file, err := os.Open(dirPath + "/" + entry.Name())
			if err != nil {
				return nil, true, fmt.Errorf("failed to open Dockerfile: %w", err)
			}
			defer file.Close()

			dockerfile, err := dockerfile.New(entry.Name(),file)
			if err != nil {
				return nil, true, fmt.Errorf("failed to parse Dockerfile: %w", err)
			}

			return dockerfile, true, nil
		}
	}

	return nil, false, nil
}