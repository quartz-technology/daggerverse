package compose

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/module/object"
	"dagger.io/dockersdk/module/proxy"
)

// allFunc is a function that starts all services using
// a proxy module to group them together.
//
// It MUST be registered AFTER all services have been registered.
// The proxy is a simple duplication of: github.com/kpenfound/dagger-modules/proxy@v0.2.5 module.
type allFunc struct {
	c *Compose
}

func (u *allFunc) up(services []*proxy.Service) *dagger.Container {
	proxy := proxy.New()

	for _, service := range services {
		proxy = proxy.WithService(service)
	}

	return proxy.Service()
}

func (u *allFunc) Invoke(ctx context.Context, state object.State, input object.InputArgs) (object.Result, error) {
	compose, err := u.c.load(state)
	if err != nil {
		return nil, fmt.Errorf("failed to load object state: %w", err)
	}

	services := []*proxy.Service{}
	for _, service := range u.c.dockercompose.Services() {
		service := &serviceFunc{c: compose, service: service}
		servicePrefix := fmt.Sprintf("%s_", service.service.Name())

		serviceInput := input
		for argName, argValue := range input {
			if strings.HasPrefix(argName, servicePrefix) {
				serviceInput[strings.TrimPrefix(argName, servicePrefix)] = argValue
			}
		}

		serviceCtr, err := service.ToService(ctx, state, serviceInput)
		if err != nil {
			return nil, fmt.Errorf("failed to get service %s: %w", service.service.Name(), err)
		}

		services = append(services, serviceCtr)
	}

	return (*allFunc).up(
		&allFunc{c: compose},
		services,
	), nil
}

func (u *allFunc) Arguments() []*object.FunctionArg {
	// This method should not be called for this function
	return nil
}

func (u *allFunc) AddTypeDefToObject(ctx context.Context, mod *dagger.Module, obj *dagger.TypeDef) (*dagger.Module, *dagger.TypeDef) {
	args := []*object.FunctionArg{}

	serviceNames := []string{}
	for name, service := range u.c.funcMap {
		serviceNames = append(serviceNames, name)

		serviceArgs := service.Arguments()
		for _, arg := range serviceArgs {
			args = append(args, &object.FunctionArg{
				// Prefix the argument name with the service name to avoid colission
				Name: fmt.Sprintf("%s_%s", name, arg.Name),
				Type: arg.Type,
				Opts: arg.Opts,
			})
		}
	}

	typedef := dag.Function("All", dag.TypeDef().WithObject("Container")).
		WithDescription(fmt.Sprintf("Start all service containers (%s)", strings.Join(serviceNames, ", ")))

	for _, arg := range args {
		typedef = typedef.WithArg(arg.Name, arg.Type, arg.Opts)
	}

	return mod, obj.
		WithFunction(typedef)
}
