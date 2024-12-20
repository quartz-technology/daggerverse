package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger/dag"
	"dagger.io/magicsdk/module"
)

func main() {
	ctx := context.Background()
	defer dag.Close()

	mod, err := module.Build("Test", "/app")
	if err != nil {
		panic(fmt.Errorf("failed to build module: %w", err))
	}

	if err := mod.Dispatch(ctx); err != nil {
		os.Exit(2)
	}
}