package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/codebase"
)

func main() {
	ctx := context.Background()
	defer dag.Close()

	name, err := dag.CurrentModule().Name(ctx)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to get user's codebase: %w", err))

		os.Exit(2)
	}

	formattedName := strings.ToUpper(string(name[0])) + name[1:]

	codebase, err := codebase.New(ctx)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to get user's codebase: %w", err))

		os.Exit(2)
	}

	if err := codebase.ToModule(formattedName).Dispatch(ctx); err != nil {
		os.Exit(2)
	}
}
