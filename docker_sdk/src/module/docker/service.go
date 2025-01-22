package docker

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/codebase/dockercompose"
	"dagger.io/dockersdk/module/object"
	"dagger.io/dockersdk/utils"
)

type serviceFunc struct {
	service *dockercompose.Service
}

func (s *serviceFunc) container(image string) *dagger.Container {	
	ctr := dag.Container().From(image)

	if workdir := s.service.Workdir(); workdir != "" {
		ctr = ctr.WithWorkdir(workdir)
	}

	return ctr
}

func (s *serviceFunc) Invoke(ctx context.Context, _ object.State, input object.InputArgs) (object.Result, error) {
	image := utils.LoadArgument[string]("image", input)

	return s.container(image), nil
}

func (s *serviceFunc) AddTypeDefToObject(ctx context.Context, mod *dagger.Module, object *dagger.TypeDef) (*dagger.Module, *dagger.TypeDef) {
	typedef := dag.
		Function(s.service.Name(), dag.TypeDef().WithObject("Container")).
		WithDescription(fmt.Sprintf("Create a %s service container", s.service.Name())).
		WithArg("image", dag.TypeDef().WithKind(dagger.TypeDefKindStringKind), dagger.FunctionWithArgOpts{
			DefaultValue: utils.LoadDefaultValue(s.service.Image()),
			Description: "Image to use for the service",
		})

	return mod, object.WithFunction(typedef)
}
