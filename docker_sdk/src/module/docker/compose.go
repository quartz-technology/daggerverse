package docker

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/module/compose"
	"dagger.io/dockersdk/module/object"
)

type composeFunc struct {
	d *Docker
}

func (c *composeFunc) compose() *compose.Compose {
	return compose.New(c.d.Dir, c.d.dockercomposeFile)
}

func (c *composeFunc) Invoke(ctx context.Context, state object.State, input object.InputArgs) (object.Result, error) {
	if c.d.dockercomposeFile == nil {
		return nil, fmt.Errorf("docker-compose file not loaded")
	}
	
	docker, err := c.d.load(state)
	if err != nil {
		return nil, fmt.Errorf("failed to object state: %w", err)
	}

	return (*composeFunc).compose(&composeFunc{d: docker}), nil
}

func (c *composeFunc) Arguments() []*object.FunctionArg {
	// This method should never be called for this function
	return nil
}

func (c *composeFunc) AddTypeDefToObject(ctx context.Context, mod *dagger.Module, object *dagger.TypeDef) (*dagger.Module, *dagger.TypeDef) {
	typedef := dag.Function("Compose", dag.TypeDef().WithObject("Compose")).
		WithDescription("Manage docker compose services")

	return mod, object.WithFunction(typedef)
}