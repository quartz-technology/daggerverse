package compose

import (
	"context"
	"encoding/json"
	"fmt"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/codebase/dockercompose"
	"dagger.io/dockersdk/module/object"
	"dagger.io/dockersdk/utils"
)

type Compose struct {
	Dir *dagger.Directory

	dockercompose *dockercompose.DockerCompose
	funcMap       map[string]object.Function
}

func New(
	dir *dagger.Directory,
	dockercomposeFile *dockercompose.DockerCompose,
) *Compose {
	c := &Compose{
		Dir:           dir,
		dockercompose: dockercomposeFile,
		funcMap:       make(map[string]object.Function),
	}

	for _, service := range dockercomposeFile.Services() {
		c.funcMap[service.Name()] = &serviceFunc{c: c, service: service}
	}

	// Add up for all services
	c.funcMap["All"] = &allFunc{c: c}

	return c
}

func (c *Compose) Name() string {
	return "Compose"
}

func (c *Compose) Description() string {
	return "Manage compos services"
}

func (c *Compose) New(input object.InputArgs) object.Object {
	var dir *dagger.Directory

	if input["dir"] != nil {
		dir = utils.LoadDirectoryFromID([]byte(input["dir"]))
	}

	return &Compose{
		Dir:           dir,
		dockercompose: c.dockercompose,
	}
}

func (c *Compose) AddTypeDef(ctx context.Context) dagger.WithModuleFunc {
	return func(mod *dagger.Module) *dagger.Module {
		object := dag.TypeDef().WithObject(c.Name())

		for _, fct := range c.funcMap {
			mod, object = fct.AddTypeDefToObject(ctx, mod, object)
		}

		return mod.WithObject(object)
	}
}

func (c *Compose) Load(state object.State) (object.Object, error) {
	return c.load(state)
}

func (c *Compose) load(state object.State) (*Compose, error) {
	parentMap := make(map[string]interface{})
	err := json.Unmarshal(state, &parentMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal parent object: %w", err)
	}

	cpyCompose := &Compose{
		dockercompose: c.dockercompose,
		funcMap:       c.funcMap,
	}

	if parentMap["Dir"] != nil {
		cpyCompose.Dir = dag.LoadDirectoryFromID(dagger.DirectoryID(parentMap["Dir"].(string)))
	}

	return cpyCompose, nil
}

func (c *Compose) Invoke(ctx context.Context, state object.State, fnName string, input object.InputArgs) (object.Result, error) {
	if c.funcMap[fnName] == nil {
		return nil, fmt.Errorf("unknown function %s", fnName)
	}

	return c.funcMap[fnName].Invoke(ctx, state, input)
}
