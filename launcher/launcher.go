package main

import (
	"context"
	"time"
)

type Launcher struct{}

// Publish adds the module to the daggerverse.
func (l *Launcher) Publish(ctx context.Context, module *Directory, path Optional[string]) (string, error) {
	modulePath := path.GetOr(".")

	return dag.Dagger().Cli("0.9.2").
		Ctr().
		WithMountedDirectory("/module", module).
		WithWorkdir("/module").
		WithEnvVariable("CACHE_BURST", time.DateTime).
		WithExec(
			[]string{"mod", "publish", "-m", modulePath, "-f"},
			ContainerWithExecOpts{ExperimentalPrivilegedNesting: true},
		).
		Stdout(ctx)
}
