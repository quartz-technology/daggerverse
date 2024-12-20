package integration

import (
	"context"
	"encoding/json"
	"fmt"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/magicsdk/codebase"
	"dagger.io/magicsdk/invocation"
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

func (d *Docker) Exist() bool {
	return d.supported
}

func (d *Docker) TypeDef() *dagger.TypeDef {
	return dag.
		TypeDef().
		WithObject("Docker").
		WithFunction(
			dag.Function("Build", dag.TypeDef().WithObject("Container")).
				WithDescription("Build a container from a Dockerfile"),
		)
}

func (d *Docker) New(dir *dagger.Directory) Integration {
	return &Docker{
		Dir: dir,
	}
}

func (d *Docker) Build(ctx context.Context) (*dagger.Container, error) {
	return dag.Container().From("alpine").WithDirectory("/app", d.Dir), nil
}

func (d *Docker) Invoke(ctx context.Context, invocation *invocation.Invocation) (_ any, err error) {
	switch invocation.FnName {
	case "Build":
		{
			var parent Docker
			err = json.Unmarshal(invocation.ParentJSON, &parent)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal parent object: %w", err)
			}

			fmt.Println(string(invocation.ParentJSON))

			return (*Docker).Build(&parent, ctx)
		}
	default: 
		return nil, fmt.Errorf("unknown function %s", invocation.FnName)
	}
}
