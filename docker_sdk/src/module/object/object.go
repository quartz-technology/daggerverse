package object

import (
	"context"

	"dagger.io/dagger"
)

// State represents the state of an object.
type State []byte

// InputArgs is a map of string to byte slice, used for function inputs.
type InputArgs map[string][]byte

// Result denotes the result of a function invocation.
type Result any

// FunctionArg defines an argument for a function including its name, type,
// and options.
type FunctionArg struct {
	// Name is the name of the function argument.
	Name string

	// Type specifies the type definition of the argument.
	Type *dagger.TypeDef

	// Opts provides additional options for the function argument.
	Opts dagger.FunctionWithArgOpts
}

// Function represents an interface for functions.
type Function interface {
	// AddTypeDefToObject adds a type definition to the specified module.
	//
	// It takes as argument the DockerSDK module and the function's object TypeDef
	// and returns them with the updated module/object.
	AddTypeDefToObject(context.Context, *dagger.Module, *dagger.TypeDef) (*dagger.Module, *dagger.TypeDef)

	// Invoke calls the function with the provided state and input arguments,
	// returning the result or an error.
	Invoke(ctx context.Context, state State, input InputArgs) (Result, error)

	// Arguments returns a slice of the function arguments.
	Arguments() []*FunctionArg
}

// Object interface defines methods for handling module-related objects.
type Object interface {
	// Name returns the name of the object.
	Name() string

	// Description provides a description of the object.
	Description() string

	// AddTypeDef returns a function to exec to add the object's type definition 
	// to the module.
	AddTypeDef(context.Context) dagger.WithModuleFunc

	// Load reconstruct an object from the given state.
	Load(state State) (Object, error)

	// New creates a new object with the provided input arguments.
	New(input InputArgs) Object

	// Invoke executes a function by its name on the object, using the specified
	// state and input, and returns the result or an error.
	Invoke(ctx context.Context, state State, fnName string, input InputArgs) (Result, error)
}
