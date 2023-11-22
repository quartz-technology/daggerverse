package main

import (
	"context"
	"fmt"
)

var DefaultVersion = "v1.54-alpine"

type GolangciLint struct {
	Config *File
}

// WithVersion returns a Container  configured to use golangci-lint.
func (c *GolangciLint) WithVersion(version string) *Container {
	return dag.
		Container().
		From(fmt.Sprintf("golangci/golangci-lint:%s", version)).
		WithEntrypoint([]string{"golangci-lint"})
}

func (c *GolangciLint) WithConfig(file *File) *GolangciLint {
	c.Config = file

	return c
}

func (c *GolangciLint) Lint(ctx context.Context, directory *Directory) (string, error) {
	ctr := c.WithVersion(DefaultVersion)

	ctr = dag.
		Golang().
		WithContainer(ctr).
		WithSource(directory).
		Container()

	if c.Config != nil {
		ctr = ctr.WithMountedFile("/src/.golangci.yml", c.Config)
	}

	return ctr.
		WithExec([]string{"run", "-v", "./..."}).
		Stdout(ctx)
}
