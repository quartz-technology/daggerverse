package object

import (
	"context"

	"dagger.io/dagger"
)

type State []byte
type InputArgs map[string][]byte
type Result any

type FunctionArg struct {
	Name string
	Type *dagger.TypeDef
	Opts dagger.FunctionWithArgOpts
}

type Function interface  {
	AddTypeDefToObject(context.Context, *dagger.Module, *dagger.TypeDef) (*dagger.Module, *dagger.TypeDef)
	Invoke(ctx context.Context, state State,input InputArgs) (Result, error)
	Arguments() ([]*FunctionArg)
}

type Object interface {
	Name() string

	Description() string

	AddTypeDef(context.Context) dagger.WithModuleFunc

	Load(state State) (Object, error)

	New(input InputArgs) Object

	Invoke(ctx context.Context, state State, fnName string, input InputArgs) (Result, error)
}