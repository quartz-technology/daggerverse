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

// up creates a proxy module with the given services.
func (u *allFunc) up(services []*proxy.Service) *dagger.Container {
	proxy := proxy.New()

	for _, service := range services {
		proxy = proxy.WithService(service)
	}

	return proxy.Service()
}

// Invoke executes the "All" function with the given state and input arguments.
func (u *allFunc) Invoke(ctx context.Context, state object.State, input object.InputArgs) (object.Result, error) {
	compose, err := u.c.load(state)
	if err != nil {
		return nil, fmt.Errorf("failed to load object state: %w", err)
	}

	services := []*proxy.Service{}
	for _, service := range u.c.dockercompose.Services() {
		if u.c.runningServices[service.Name()] != nil {
			fmt.Printf("service %s is already running ; exposing it to the proxy\n", service.Name())

			services = append(services, u.c.runningServices[service.Name()])

			continue
		}

		fmt.Printf("service %s is not running yet; starting it\n", service.Name())

		service := &serviceFunc{c: compose, service: service, asDep: true}

		serviceCtr, err := service.ToService(ctx, state, input)
		if err != nil {
			return nil, fmt.Errorf("failed to get service %s: %w", service.service.Name(), err)
		}

		services = append(services, serviceCtr)
		u.c.runningServices[service.service.Name()] = serviceCtr
	}

	return (*allFunc).up(
		&allFunc{c: compose},
		services,
	), nil
}

// Arguments is a placeholder method not invoked for this function
// required to implements the object.Function interface.
//
// This function should never be called for this function.
func (u *allFunc) Arguments() []*object.FunctionArg {
	return nil
}

// AddTypeDefToObject adds "All" function definition to the given Dagger module's object.
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
