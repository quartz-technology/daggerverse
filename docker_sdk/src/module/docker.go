package module

import (
	"context"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/module/docker"
	"dagger.io/dockersdk/module/object"
)

type dockerFunc struct {
	d *docker.Docker
}

func (d *dockerFunc) Invoke(_ context.Context, _ object.State, input object.InputArgs) (object.Result, error) {
	return d.d.New(input), nil
}

func (d *dockerFunc) Arguments() []*object.FunctionArg {
	// This method should never be called for this function
	return nil
}

func (d *dockerFunc) AddTypeDefToObject(ctx context.Context, mod *dagger.Module, object *dagger.TypeDef) (*dagger.Module, *dagger.TypeDef) {
	return mod, object.WithFunction(
		dag.Function(d.d.Name(), dag.TypeDef().WithObject(d.d.Name())).
			WithDescription(d.d.Description()).
			WithArg("dir", dag.TypeDef().WithObject("Directory").WithOptional(true), dagger.FunctionWithArgOpts{
				DefaultPath: ".",
			}))
}
