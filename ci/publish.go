package main

import (
	"context"
)

func (c *Ci) Publish(ctx context.Context) error {
	//modules := source()

	//	directories := []string{"golang", "golangci-lint", "node", "redis", "postgres", "launcher"}

	//	eg, gctx := errgroup.WithContext(ctx)

	_, err := dag.Launcher().Publish(ctx, source(), LauncherPublishOpts{Path: "launcher"})

	return err
}
