package object

import (
	"context"

	"dagger.io/dagger"
)

type State []byte
type InputArgs map[string][]byte
type Result any

type Object interface {
	Name() string

	Description() string

	AddTypeDef(context.Context) dagger.WithModuleFunc

	Load(state State) (Object, error)

	New(input InputArgs) Object

	Invoke(ctx context.Context, fnName string, input InputArgs) (Result, error)
}