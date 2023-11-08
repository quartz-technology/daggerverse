package main

import (
	"context"
	"fmt"
	"time"
)

type Launcher struct{}

// dagger returns a dagger container with the specified dagger version.
// If no version is specified, the latest version will be used.
func dagger(version string) *Container {
	version = fmt.Sprintf("DAGGER_VERSION=%s", version)

	return dag.Container().From("alpine:latest").
		WithExec([]string{"apk", "add", "curl"}).
		WithExec([]string{"sh", "-c", fmt.Sprintf("curl -L https://dl.dagger.io/dagger/install.sh | %s sh", version)}).
		WithEntrypoint([]string{"/bin/dagger"})
}

func (l *Launcher) Publish(ctx context.Context, module *Directory, path Optional[string]) (string, error) {
	modulePath := path.GetOr(".")

	return dagger("0.9.2").
		WithMountedDirectory("/module", module).
		WithWorkdir("/module").
		WithEnvVariable("CACHE_BURST", time.DateTime).
		WithExec(
			[]string{"mod", "publish", "-m", modulePath, "-f"},
			ContainerWithExecOpts{ExperimentalPrivilegedNesting: true},
		).
		Stdout(ctx)
}
