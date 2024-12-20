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
)

type Docker struct {
	Dir *dagger.Directory

	supported bool
}

func DockerIntegration(code *codebase.Codebase) (Integration, error) {
	_, err, exist := code.LookupFile("Dockerfile")
	if err != nil {
		return nil, fmt.Errorf("failed to lookup for Dockerfile: %w", err)
	}

	return &Docker{
		supported: exist,
	}, nil
}

func (d *Docker) Description() string {
	return "Access docker functions"
}

func (d *Docker) Exist() bool {
	return d.supported
}

func (d *Docker) TypeDef() *dagger.TypeDef {
	return dag.
		TypeDef().
		WithObject("Docker").
		WithFunction(
			dag.Function("Build", dag.TypeDef().WithObject("Container")).
				WithDescription("Build a container the Dockerfile present in the app"),
		)
}

func (d *Docker) New(invocation *invocation.Invocation) Integration {
	// Workaround to parse argument since `UnmarshalJSON` isn't generated for Dagger type in
	// the client library.
	// This should be fixed later to integrate a real MagicSDK
	var dir *dagger.Directory

	if invocation.InputArgs["dir"] != nil {
		dir = utils.LoadDirectoryFromID([]byte(invocation.InputArgs["dir"]))
	}

	return &Docker{
		Dir: dir,
	}
}

func (d *Docker) Build(ctx context.Context) (*dagger.Container, error) {
	return d.Dir.DockerBuild(), nil
}

func (d *Docker) Invoke(ctx context.Context, invocation *invocation.Invocation) (_ any, err error) {
	switch invocation.FnName {
	case "Build":
		// Workaround to parse argument since `UnmarshalJSON` isn't generated for Dagger type in
		// the client library.
		// This should be fixed later to integrate a real MagicSDK
		var parentMap map[string]interface{}
		err = json.Unmarshal(invocation.ParentJSON, &parentMap)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal parent object: %w", err)
		}

		parent := Docker{
			Dir: dag.LoadDirectoryFromID(dagger.DirectoryID(parentMap["Dir"].(string))),
		}

		return (*Docker).Build(&parent, ctx)
	default:
		return nil, fmt.Errorf("unknown function %s", invocation.FnName)
	}
}
