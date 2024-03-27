// Magicenv provides utility functions to inject environment from various sources.
//
// Most of modern application use environment to configure runtime variable, such as
// database connection, API key, etc.
// Usually, they are store in a file, such as .env, .env.local, .env.development, or .envrc.
//
// This package provides a way to inject those environment into a container so you do not have
// to redifine all your variable when you "daggerize" your application.

package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-envparse"
)

type Magicenv struct{}

// Load environment from a .env or .envrc type of tile and inject it into a container.
//
// All variables are loaded as plaintext so be careful with sensitive secrets.
func (m *Magicenv) LoadEnv(
	ctx context.Context,

	// The container to inject the environment in.
	ctr *Container,
	// The path to the environment file to load.
	//
	//+optional
	//+default=".env"
	path string,
) (*Container, error) {
	// Get the content of the .env file
	envFile, err := ctr.File(path).Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	// Parse it to get a map
	envMap, err := envparse.Parse(strings.NewReader(envFile))
	if err != nil {
		return nil, fmt.Errorf("failed to parse environment file: %w", err)
	}

	// Inject all the variables as secret
	for key, value := range envMap {
		ctr = ctr.WithEnvVariable(key, value, ContainerWithEnvVariableOpts{ Expand: true })
	}

	return ctr, nil
}
