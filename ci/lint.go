package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
)

func (c *Ci) Lint(ctx context.Context) error {
	lintCtr := dag.
		GolangciLint().
		WithVersion("v1.54-alpine")

	ctr := dag.
		Golang().
		WithContainer(lintCtr).
		WithSource(source())

	eg, gctx := errgroup.WithContext(ctx)

	// Execute linter on given files.
	lintFct := func(files string) func() error {
		return func() error {
			out, err := ctr.Exec([]string{"run", "-v", files}).Stdout(gctx)
			if err != nil {
				return err
			}

			fmt.Println(out)

			return nil
		}
	}

	eg.Go(lintFct("golang/*.go"))
	eg.Go(lintFct("golangci-lint/*.go"))
	eg.Go(lintFct("node/*.go"))

	return eg.Wait()
}
