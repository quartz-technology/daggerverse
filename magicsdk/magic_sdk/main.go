package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"dagger.io/dagger/dag"
	"dagger.io/magicsdk/module"
)

func main() {
	ctx := context.Background()
	defer dag.Close()

	name, err := dag.CurrentModule().Name(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to get current module name: %w", err))
	}

	formattedName := strings.ToUpper(string(name[0])) + name[1:]

	mod, err := module.Build(formattedName, "/app")
	if err != nil {
		panic(fmt.Errorf("failed to build module: %w", err))
	}

	if err := mod.Dispatch(ctx); err != nil {
		os.Exit(2)
	}
}