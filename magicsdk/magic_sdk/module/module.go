package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/magicsdk/codebase"
	"dagger.io/magicsdk/integration"
	"dagger.io/magicsdk/invocation"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Module struct {
	name string
	codebase *codebase.Codebase
	integrations integration.Integrations
}

func Build(name string, path string) (*Module, error) {
	codebase, err := codebase.New(path)
	if err != nil {
		return nil, fmt.Errorf("failed to build module: %w", err)
	}

	integrations, err := integration.LoadIntegrations(codebase)
	if err != nil {
		return nil, fmt.Errorf("failed to load integrations: %w", err)
	}

	return &Module{
		name: name,
		codebase: codebase,
		integrations: integrations,
	}, nil
}

func (m *Module) TypeDef() *dagger.Module {
	mod := dag.Module()

	mainObject := dag.TypeDef().WithObject(m.name)
	for name, integration := range m.integrations {
		mainObject = mainObject.WithFunction(
			dag.Function(name, dag.TypeDef().WithObject(name)).
				WithArg("dir", dag.TypeDef().WithObject("Directory"), dagger.FunctionWithArgOpts{
					DefaultPath: "/",
				}),
		)

		mod = mod.WithObject(integration.TypeDef())
	}

	mod = mod.WithObject(mainObject)

	return mod
}

func unwrapError(rerr error) string {
	var gqlErr *gqlerror.Error
	if errors.As(rerr, &gqlErr) {
		return gqlErr.Message
	}
	return rerr.Error()
}

func (m *Module) Dispatch(ctx context.Context) (rerr error) {
	fnCall := dag.CurrentFunctionCall()
	defer func() {
		if rerr != nil {
			if err := fnCall.ReturnError(ctx, dag.Error(unwrapError(rerr))); err != nil {
				fmt.Println("failed to return error:", err)
			}
		}
	}()

	parentName, err := fnCall.ParentName(ctx)
	if err != nil {
		return fmt.Errorf("get parent name: %w", err)
	}
	fnName, err := fnCall.Name(ctx)
	if err != nil {
		return fmt.Errorf("get fn name: %w", err)
	}
	parentJson, err := fnCall.Parent(ctx)
	if err != nil {
		return fmt.Errorf("get fn parent: %w", err)
	}
	fnArgs, err := fnCall.InputArgs(ctx)
	if err != nil {
		return fmt.Errorf("get fn args: %w", err)
	}

	inputArgs := map[string][]byte{}
	for _, fnArg := range fnArgs {
		argName, err := fnArg.Name(ctx)
		if err != nil {
			return fmt.Errorf("get fn arg name: %w", err)
		}
		argValue, err := fnArg.Value(ctx)
		if err != nil {
			return fmt.Errorf("get fn arg value: %w", err)
		}
		inputArgs[argName] = []byte(argValue)
	}

	result, err := m.Invoke(ctx, &invocation.Invocation{
		ParentJSON: []byte(parentJson),
		ParentName: parentName,
		FnName:     fnName,
		InputArgs:  inputArgs,
	})
	if err != nil {
		var exec *dagger.ExecError
		if errors.As(err, &exec) {
			return exec.Unwrap()
		}
		return err
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if err := fnCall.ReturnValue(ctx, dagger.JSON(resultBytes)); err != nil {
		return fmt.Errorf("store return value: %w", err)
	}

	return nil
}

func (m *Module) Invoke(ctx context.Context, invocation *invocation.Invocation) (_ any, err error) {
	switch invocation.ParentName {
	case m.name:
		switch invocation.FnName {
		case "Docker":
			var dir *dagger.Directory
			if invocation.InputArgs["app"] != nil {
				err = json.Unmarshal([]byte(invocation.InputArgs["app"]), &dir)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg dir", err))
				}
			}

			return m.integrations["Docker"].New(dir), nil
		default:
			return nil, fmt.Errorf("unknown function %s", invocation.FnName)
		}
	case "Docker":
		return m.integrations["Docker"].Invoke(ctx, invocation)
	case "":
		return m.TypeDef(), nil
	default:
		return nil, fmt.Errorf("unknown object %s", invocation.ParentName)
	}
}