package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/codebase"
	"dagger.io/dockersdk/module"
)

func main() {
	ctx := context.Background()
	defer dag.Close()

	name, err := dag.CurrentModule().Name(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to get current module name: %w", err))
	}

	formattedName := strings.ToUpper(string(name[0])) + name[1:]

	codebase, err := codebase.New()
	if err != nil {
		panic(fmt.Errorf("failed to get user's codebase: %w", err))
	}

	mod := module.Build(formattedName, codebase)

	if err := mod.Dispatch(ctx); err != nil {
		os.Exit(2)
	}
}