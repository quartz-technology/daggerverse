package main

import (
	"context"
	"golang.org/x/sync/errgroup"
)

type IntegrationTest struct{}

func (c *Ci) IntegrationTest() IntegrationTest {
	return IntegrationTest{}
}

func (i *IntegrationTest) Run(ctx context.Context) error {
	eg, gctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return i.Node(gctx)
	})

	eg.Go(func() error {
		return i.Redis(gctx)
	})

	return eg.Wait()
}
