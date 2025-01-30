package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/module/docker"
	"dagger.io/dockersdk/module/object"
	"dagger.io/dockersdk/utils"
)

// Module represents a configurable DockerSDK module.
type Module struct {
	// name specifies the name of the module.
	name string

	// funcMap holds functions available in this module by name.
	funcMap map[string]object.Function

	// objects contains the objects associated with this module.
	objects map[string]object.Object
}

// Build initializes a new Module with a given name and Docker configuration.
func Build(name string, docker *docker.Docker) *Module {
	baseObjects := map[string]object.Object{
		// The default object for the docker SDK
		"Docker": docker,
	}

	objects := utils.MergeObjectsMap(baseObjects, docker.Deps())

	return &Module{
		name: name,
		funcMap: map[string]object.Function{
			"Docker": &dockerFunc{d: docker},
		},
		objects: objects,
	}
}

// Name returns the name of the module.
func (m *Module) Name() string {
	return m.name
}

// typeDef constructs and returns the module's type definition.
func (m *Module) typeDef(ctx context.Context) (*dagger.Module, error) {
	mod := dag.Module()

	entrypointObject := dag.TypeDef().
		WithObject(m.name)

	for _, fct := range m.funcMap {
		mod, entrypointObject = fct.AddTypeDefToObject(ctx, mod, entrypointObject)

	}

	for _, obj := range m.objects {
		mod = mod.With(obj.AddTypeDef(ctx))
	}

	mod = mod.WithObject(entrypointObject)

	return mod, nil
}

// Dispatch handles the execution of a function call within the module context.
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

	result, err := m.invoke(ctx, parentName, []byte(parentJson), fnName, inputArgs)
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

// invoke carries out the invocation of a function within the module.
func (m *Module) invoke(ctx context.Context, parentName string, parentJSON object.State, fnName string, input object.InputArgs) (_ any, err error) {
	// If it's an empty parent name, that means we need to register the
	// module's type definition.
	if parentName == "" {
		return m.typeDef(ctx)
	}

	// If it's a top-level invocation, we build the called object.
	if parentName == m.name {
		return m.funcMap[fnName].Invoke(ctx, parentJSON, input)
	}

	// If it's an object invocation, we build the docker SDK and invoke the function
	object, err := m.objects[parentName].Load(parentJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to load docker SDK module: %w", err)
	}

	return object.Invoke(ctx, parentJSON, fnName, input)
}
