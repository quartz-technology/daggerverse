package module

import (
	"context"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/module/docker"
	"dagger.io/dockersdk/module/object"
)

// dockerFunc orchestrates Docker-related actions.
type dockerFunc struct {
	// d is an instance of the Docker Object.
	d *docker.Docker
}

// Invoke executes a Docker function with given input arguments.
func (d *dockerFunc) Invoke(_ context.Context, _ object.State, input object.InputArgs) (object.Result, error) {
	return d.d.New(input), nil
}

// Arguments returns the arguments required to invoke the Docker function.
//
// Note: This method should not be called and only
// exists to implement the object.Function interface.
func (d *dockerFunc) Arguments() []*object.FunctionArg {
	return nil
}

// AddTypeDefToObject enriches the object's definition with the Docker function.
func (d *dockerFunc) AddTypeDefToObject(ctx context.Context, mod *dagger.Module, object *dagger.TypeDef) (*dagger.Module, *dagger.TypeDef) {
	return mod, object.WithFunction(
		dag.Function(d.d.Name(), dag.TypeDef().WithObject(d.d.Name())).
			WithDescription(d.d.Description()).
			WithArg("dir", dag.TypeDef().WithObject("Directory").WithOptional(true), dagger.FunctionWithArgOpts{
				DefaultPath: ".",
			}))
}
