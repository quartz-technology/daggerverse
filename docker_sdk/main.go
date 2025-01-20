// A generated module for Dockersdk functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/dockersdk/internal/dagger"
	"fmt"
)

type Dockersdk struct {
	App *dagger.Directory

	RequiredPaths []string
}

func New(
	// Source file of the Docker SDK, this path should never be changed nor set.
	//+defaultPath="./src"
	app *dagger.Directory,
) *Dockersdk {
	return &Dockersdk{
		App: app,
	}
}

func (m *Dockersdk) ModuleRuntime(ctx context.Context, modSource *dagger.ModuleSource, introspectionJSON *dagger.File) (*dagger.Container, error) {
	runtimeBin := dag.Container().
		From("golang:1.23.2-alpine").
		WithDirectory("/src", m.App).
		WithWorkdir("/src").
		WithExec([]string{"go", "build", "-o", "/src/magic_sdk", "."}).
		File("/src/magic_sdk")

	modulePath, err := modSource.SourceRootSubpath(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get module path: %w", err)
	}

	sourceDir := modSource.ContextDirectory().Directory(modulePath)
	if sourceDir == nil {
		sourceDir = dag.Directory()
	}

	return dag.
		Container().
		From("golang:1.23.2-alpine").
		WithWorkdir("/runtime").
		WithFile("/runtime/magic_sdk", runtimeBin).
		WithDirectory("/app", sourceDir).
		WithEntrypoint([]string{"/runtime/magic_sdk"}), nil
}

// The Docker SDK does not generate any code.
func (m *Dockersdk) Codegen(ctx context.Context, modSource *dagger.ModuleSource, introspectionJSON *dagger.File) (*dagger.GeneratedCode, error) {
	return dag.GeneratedCode(dag.Directory()), nil
}