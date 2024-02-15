package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
)

func (c *Ci) Lint(ctx context.Context) error {
	eg, gctx := errgroup.WithContext(ctx)
	source := source()

	configFile := source.File(".golangci.yml")

	// Execute linter on given files.
	lintFct := func(directoryPath string) func() error {
		return func() error {
			out, err := dag.GolangciLint().
				WithConfig(configFile).
				Lint(gctx, source.Directory(directoryPath))

			fmt.Println(out)

			return err
		}
	}

	directories := []string{"golang", "golangci-lint", "node", "redis", "postgres", "dagger", "minio", "launcher"}
	for _, d := range directories {
		eg.Go(lintFct(d))
	}

	return eg.Wait()
}
