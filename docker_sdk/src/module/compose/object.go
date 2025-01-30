package compose

import (
	"context"
	"encoding/json"
	"fmt"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/codebase/dockercompose"
	"dagger.io/dockersdk/module/object"
	"dagger.io/dockersdk/module/proxy"
	"dagger.io/dockersdk/utils"
)

// Compose manages Docker Compose services.
type Compose struct {
	// Dir is the directory containing the Docker Compose configuration.
	Dir *dagger.Directory

	// dockercompose is the Docker Compose configuration object.
	dockercompose *dockercompose.DockerCompose

	// funcMap maps service names to their associated functions for operations.
	funcMap map[string]object.Function

	// runningServices maps service names to their running services.
	runningServices map[string]*proxy.Service
}

// New creates a new Compose instance with the given directory and docker-compose file.
func New(
	dir *dagger.Directory,
	dockercomposeFile *dockercompose.DockerCompose,
) *Compose {
	c := &Compose{
		Dir:             dir,
		dockercompose:   dockercomposeFile,
		funcMap:         make(map[string]object.Function),
		runningServices: make(map[string]*proxy.Service),
	}

	for _, service := range dockercomposeFile.Services() {
		c.funcMap[service.Name()] = &serviceFunc{c: c, service: service, asDep: false}
	}

	// Add a function to start all services.
	c.funcMap["All"] = &allFunc{c: c}

	return c
}

// Name returns the name of the object: "Compose".
func (c *Compose) Name() string {
	return "Compose"
}

// Description provides a brief description of Compose.
func (c *Compose) Description() string {
	return "Manage compos services"
}

// New creates a new Compose object instance with optional directory input.
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

// AddTypeDef adds the module type definition for this object with all
// its functions.
func (c *Compose) AddTypeDef(ctx context.Context) dagger.WithModuleFunc {
	return func(mod *dagger.Module) *dagger.Module {
		object := dag.TypeDef().WithObject(c.Name())

		for _, fct := range c.funcMap {
			mod, object = fct.AddTypeDefToObject(ctx, mod, object)
		}

		return mod.WithObject(object)
	}
}

// Load constructs a new Compose object from a saved state.
func (c *Compose) Load(state object.State) (object.Object, error) {
	return c.load(state)
}

// load reconstructs a new Compose from state data.
func (c *Compose) load(state object.State) (*Compose, error) {
	parentMap := make(map[string]interface{})
	err := json.Unmarshal(state, &parentMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal parent object: %w", err)
	}

	cpyCompose := &Compose{
		dockercompose:   c.dockercompose,
		funcMap:         c.funcMap,
		runningServices: c.runningServices,
	}

	if parentMap["Dir"] != nil {
		cpyCompose.Dir = dag.LoadDirectoryFromID(dagger.DirectoryID(parentMap["Dir"].(string)))
	}

	return cpyCompose, nil
}

// Invoke executes a function associated from its name with its object's state and input.
func (c *Compose) Invoke(ctx context.Context, state object.State, fnName string, input object.InputArgs) (object.Result, error) {
	if c.funcMap[fnName] == nil {
		return nil, fmt.Errorf("unknown function %s", fnName)
	}

	return c.funcMap[fnName].Invoke(ctx, state, input)
}
