package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
)

func (i *IntegrationTest) Node(ctx context.Context) error {
	source := daggerRepository().Directory("./sdk/nodejs")

	nodeCtr := dag.
		Node().
		WithVersion("20-alpine3.17").
		WithSource(source)

	eg, gctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return i.node(gctx, nodeCtr.WithYarn())
	})

	eg.Go(func() error {
		return i.node(gctx, nodeCtr.WithNpm())
	})

	return eg.Wait()
}

func (i *IntegrationTest) node(ctx context.Context, ctr *Node) error {
	eg, gctx := errgroup.WithContext(ctx)

	// Install dependencies
	ctr = ctr.Install([]string{})

	eg.Go(func() error {
		fmt.Println("Use Run to check version")

		_, err := ctr.Run([]string{"--version"}).Sync(gctx)

		return err
	})

	eg.Go(func() error {
		fmt.Println("Try Linter")

		out, err := ctr.Lint(gctx)
		if err != nil {
			return err
		}

		fmt.Println(out)

		return nil
	})

	eg.Go(func() error {
		fmt.Println("Try to build")

		files, err := ctr.
			Build().
			Container().
			Directory("dist").
			Entries(gctx)

		if err != nil {
			return err
		}

		if len(files) == 0 {
			return fmt.Errorf("no artifacts has been produced during build")
		}

		return nil
	})

	//	eg.Go(func() error {
	//		fmt.Println("Try to publish with dry-run")
	//
	//
	//		secretToken := dag.SetSecret("TOKEN", os.Getenv("NODE_PUBLISH_TOKEN"))
	//
	//		out, err := ctr.Publish(ctx, NodePublishOpts{
	//			DryRun:  true,
	//			Token:   secretToken,
	//			Version: "1.0.0",
	//		})
	//
	//		if err != nil {
	//			return err
	//		}
	//
	//		fmt.Println(out)
	//
	//		return nil
	//	})

	return eg.Wait()
}
