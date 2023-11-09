package main

import (
	"fmt"
	"strings"
)

type Dagger struct{}

// CLI wraps dagger binary into a specialized container that set options to use
// Dagger inside dagger.
// Each supported command are set with `ExperimentalPrivilegedNesting` option
// to ensure dagger can work inside dagger.
type CLI struct {
	Ctr *Container
}

// Install returns the same Container with the Dagger CLI installed in it.
// This container must have `curl` installed to fetch the CLI.
//
// This can be used to provide a container with source code already installed
// in it.
func (d *Dagger) Install(container *Container, version string) *CLI {
	// Remove v if it prefixes the version
	version = strings.TrimPrefix(version, "v")

	// Format with the dagger version as environment variable
	version = fmt.Sprintf("DAGGER_VERSION=%s", version)

	ctr := container.
		WithExec([]string{"sh", "-c", fmt.Sprintf("curl -L https://dl.dagger.io/dagger/install.sh | %s sh", version)}).
		WithEntrypoint([]string{"/bin/dagger"})

	return &CLI{
		Ctr: ctr,
	}
}

// CLI returns a ready to use Dagger container with CLI installed.
func (d *Dagger) CLI(version string) *CLI {
	ctr := dag.
		Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "curl"})

	return d.Install(ctr, version)
}

// Container returns the CLI's Container.
func (c *CLI) Container() *Container {
	return c.Ctr
}
