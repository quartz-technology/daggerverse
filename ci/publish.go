package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"strings"
)

func (c *Ci) Publish(ctx context.Context) (string, error) {
	modules := repo()

	directories := []string{"golang", "golangci-lint", "node", "redis", "postgres", "launcher"}

	eg, gctx := errgroup.WithContext(ctx)

	var ref []string

	for _, d := range directories {
		eg.Go(func() error {
			out, err := dag.Launcher().Publish(gctx, modules, LauncherPublishOpts{Path: d})

			ref = append(ref, out)
			return err
		})
	}

	err := eg.Wait()
	if err != nil {
		return "failed to publish modules", err
	}

	return fmt.Sprintf("modules published:\n %s", strings.Join(ref, "\n")), nil
}
