package main

import (
	"context"
)

func (c *Ci) Publish(ctx context.Context) (string, error) {
	//modules := source()

	//	directories := []string{"golang", "golangci-lint", "node", "redis", "postgres", "launcher"}

	//	eg, gctx := errgroup.WithContext(ctx)

	return dag.Launcher().Publish(ctx, repo(), LauncherPublishOpts{Path: "launcher"})
}
