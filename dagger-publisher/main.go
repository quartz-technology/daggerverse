package main

import (
	"context"
	"fmt"
	"path/filepath"
)

type DaggerPublisher struct {
	Container *Container
}

const DefaultDaggerVersion = "0.9.10"

func New(
	// +optional
	version string,
) *DaggerPublisher {
	if version == "" {
		version = DefaultDaggerVersion
	}

	return &DaggerPublisher{
		Container: dag.
			Container().
			From("alpine:3.19.1").
			WithExec([]string{"apk", "add", "curl"}).
			WithExec([]string{
				"sh", "-c", 
				fmt.Sprintf("curl -L https://dl.dagger.io/dagger/install.sh | %s sh", fmt.Sprintf("DAGGER_VERSION=%s", version)),
			}).
			WithEntrypoint([]string{"/bin/dagger"}),
	}
}

// Publish executes the publish command to upload the module
// to the Daggerverse.
func (d *DaggerPublisher) Publish(
	ctx context.Context,

	// The repository to use the Dagger CLI on
	repository *Directory,

	// The path to the module to publish
	// +optional
	path string,
) (string, error) {
	return d.
		Container.
		WithWorkdir("/module").
		WithDirectory("/module", repository).
		WithWorkdir(filepath.Join("/module", path)).
		WithExec([]string{"publish"},
			ContainerWithExecOpts{ExperimentalPrivilegedNesting: true},
		).
		Stdout(ctx)
}