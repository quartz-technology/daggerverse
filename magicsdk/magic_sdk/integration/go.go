package integration

import (
	"context"
	"encoding/json"
	"fmt"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/magicsdk/codebase"
	"dagger.io/magicsdk/invocation"
	"dagger.io/magicsdk/utils"
	"golang.org/x/mod/modfile"
)

type Go struct {
	Dir *dagger.Directory

	version   string
	supported bool
}

func GoIntegration(code *codebase.Codebase) (Integration, error) {
	gomod, err, exist := code.LookupFile("go.mod")
	if err != nil {
		return nil, fmt.Errorf("failed to lookup go.mod: %w", err)
	}

	if !exist {
		return &Go{supported: false}, nil
	}

	stat, err := gomod.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get stat on go.mod file")
	}

	gomodContent := make([]byte, stat.Size())
	_, err = gomod.Read(gomodContent)
	if err != nil {
		return nil, fmt.Errorf("failed to read go.mod: %w", err)
	}

	modfile, err := modfile.Parse("go.mod", gomodContent, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse go.mod: %w", err)
	}

	return &Go{
		version:   modfile.Go.Version,
		supported: true,
	}, nil
}

func (g *Go) Description() string {
	return fmt.Sprintf("Access function to manage your Go project (version %s)", g.version)
}

func (g *Go) Exist() bool {
	return g.supported
}

func (g *Go) TypeDef() *dagger.TypeDef {
	return dag.
		TypeDef().
		WithObject("Go").
		WithFunction(
			dag.Function("Container", dag.TypeDef().WithObject("Container")).
				WithDescription("Create a Golang development container for your project"),
		)
}

func (g *Go) New(invocation *invocation.Invocation) Integration {
	// Workaround to parse argument since `UnmarshalJSON` isn't generated for Dagger type in
	// the client library.
	// This should be fixed later to integrate a real MagicSDK
	var dir *dagger.Directory

	if invocation.InputArgs["dir"] != nil {
		dir = utils.LoadDirectoryFromID([]byte(invocation.InputArgs["dir"]))
	}

	return &Go{
		Dir: dir,
	}
}

func (g *Go) Container() (*dagger.Container, error) {
	return dag.
		Container().
		From(fmt.Sprintf("golang:%s-alpine", g.version)).
		WithDirectory("/src", g.Dir).
		WithWorkdir("/src"), nil
}

func (g *Go) Invoke(ctx context.Context, invocation *invocation.Invocation) (_ any, err error) {
	switch invocation.FnName {
	case "Container":
		// Workaround to parse argument since `UnmarshalJSON` isn't generated for Dagger type in
		// the client library.
		// This should be fixed later to integrate a real MagicSDK
		var parentMap map[string]interface{}
		err = json.Unmarshal(invocation.ParentJSON, &parentMap)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal parent object: %w", err)
		}

		parent := Go{
			Dir: dag.LoadDirectoryFromID(dagger.DirectoryID(parentMap["Dir"].(string))),
			version: g.version,
		}

		return (*Go).Container(&parent)
	default:
		return nil, fmt.Errorf("unknown function %s", invocation.FnName)
	}
}
