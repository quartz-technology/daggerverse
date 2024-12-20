package main

import (
	"context"
	"dagger/magicsdk/internal/dagger"
)

type Magicsdk struct {
	App *dagger.Directory

	RequiredPaths []string
}

func New(
	//+defaultPath="./magic_sdk"
	app *dagger.Directory,
) *Magicsdk {
	return &Magicsdk{
		App: app,
	}
}

func (m *Magicsdk) ModuleRuntime(ctx context.Context, modSource *dagger.ModuleSource, introspectionJSON *dagger.File) (*dagger.Container, error) {
	runtimeBin := dag.Container().
		From("golang:1.23.2-alpine").
		WithDirectory("/src", m.App).
		WithWorkdir("/src").
		WithExec([]string{"go", "build", "-o", "/src/magic_sdk", "."}).
		File("/src/magic_sdk")

	return dag.
		Container().
		From("golang:1.23.2-alpine").
		WithWorkdir("/runtime").
		WithFile("/runtime/magic_sdk", runtimeBin).
		WithDirectory("/app", modSource.ContextDirectory()).
		WithEntrypoint([]string{"/runtime/magic_sdk"}), nil
}

// MagicSDK doesn't have any codegen logic
func (m *Magicsdk) Codegen(ctx context.Context, modSource *dagger.ModuleSource, introspectionJSON *dagger.File) (*dagger.GeneratedCode, error) {
	return dag.GeneratedCode(dag.Directory()), nil
}
