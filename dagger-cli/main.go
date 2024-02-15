package main

import (
	"context"
	"fmt"
)

type DaggerCli struct {
	Container *Container
}

const DefaultVersion = "v0.9.10"

func New(
	// The module to use the Dagger CLI with
	module *Directory,

	// +optional
	version string,
) *DaggerCli {
	if version == "" {
		version = DefaultVersion
	}

	return &DaggerCli{
		Container: dag.
			Container().
			From("alpine:3.19.1").
			WithExec([]string{"apk", "add", "curl"}).
			WithExec([]string{
				"sh", "-c", 
				fmt.Sprintf("curl -L https://dl.dagger.io/dagger/install.sh | %s sh", fmt.Sprintf("DAGGER_VERSION=%s", version)),
			}).
			WithWorkdir("/module").
			WithDirectory("/module", module).
			WithEntrypoint([]string{"/bin/dagger"}),
	}
}

// Publish executes the publish command to upload the module
// to the Daggerverse.
func (d *DaggerCli) Publish(
	ctx context.Context,
) (string, error) {
	return d.
		Container.
		WithExec([]string{"publish"},
			ContainerWithExecOpts{ExperimentalPrivilegedNesting: true},
		).
		Stdout(ctx)
}

// Call executes the call command and returns its result.
func (d *DaggerCli) Call(
	ctx context.Context,

	args ...string,
) (string, error) {
	args = append([]string{"call"}, args...)

	return d.
		Container.
		WithExec(
			args,
			ContainerWithExecOpts{ExperimentalPrivilegedNesting: true}).
		Stdout(ctx)
}