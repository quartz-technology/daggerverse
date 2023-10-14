package main

import "fmt"

type GolangciLint struct{}

// WithVersion returns a Container  configured to use golangci-lint.
func (c *GolangciLint) WithVersion(version string) *Container {
	return dag.
		Container().
		From(fmt.Sprintf("golangci/golangci-lint:%s", version)).
		WithEntrypoint([]string{"golangci-lint"})
}
