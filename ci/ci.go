package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
)

type Ci struct{}

func source() *Directory {
	return dag.Host().Directory(".", HostDirectoryOpts{
		Include: []string{"**/*.go", "**.*go.mod", "**/*go.sum"},
	})
}

func (m *Ci) Lint(ctx context.Context) error {
	lintCtr := dag.
		GolangciLint().
		WithVersion("v1.54-alpine")

	ctr := dag.
		Golang().
		WithContainer(lintCtr).
		WithSource(source())

	eg, gctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		out, err := ctr.Exec([]string{"run", "-v", "golang/*.go"}).Stdout(gctx)
		if err != nil {
			return err
		}

		fmt.Println(out)

		return nil
	})

	eg.Go(func() error {
		out, err := ctr.Exec([]string{"run", "-v", "golangci-lint/*.go"}).Stdout(gctx)
		if err != nil {
			return err
		}

		fmt.Println(out)

		return nil
	})

	eg.Go(func() error {
		out, err := ctr.Exec([]string{"run", "-v", "node/*.go"}).Stdout(gctx)
		if err != nil {
			return err
		}

		fmt.Println(out)

		return nil
	})

	return nil
}
