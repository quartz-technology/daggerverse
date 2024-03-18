package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

type introspectionArgs struct {
	Name string
}

type introspectionFunction struct {
	Name string
	Args []introspectionArgs
}

type introspectionObject struct {
	Name        string
	Constructor introspectionFunction
	Functions   []introspectionFunction
}

// IntrospectModule returns an structured representation of objects composing a module.
func (d *DaggerverseCockpit) introspectModule(
	ctx context.Context,
	module *Directory,
) ([]introspectionObject, error) {
	introspectionResult, err := d.
		CLI("10.0.2").
		Container.
		WithWorkdir("/app").
		WithMountedDirectory("/app", module).
		WithNewFile("/app/introspection.graphql", ContainerWithNewFileOpts{
			Contents: introspectionQuery,
		}).
		WithExec([]string{"query", "--doc", "introspection.graphql"}, ContainerWithExecOpts{ExperimentalPrivilegedNesting: true}).
		Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the introspection query: %w", err)
	}

	result := gjson.Get(introspectionResult, "host.directory.asModule.initialize.objects").Array()

	objects := make([]introspectionObject, len(result))
	for i, object := range result {
		if err := json.Unmarshal([]byte(object.Get("asObject").String()), &objects[i]); err != nil {
			return nil, fmt.Errorf("could not unmarshal the module object: %w", err)
		}
	}

	return objects, nil
}

// introspectionQuery is a Dagger GraphQL query used to
// introspect a module.
var introspectionQuery = `
	query {
		host {
		  directory(path: ".") {
			asModule {
			  initialize {
				description
				objects {
				  asObject {
					name
					description
					constructor {
					  description
					  args {
						name
						description
					  }
					}
					functions {
					  name
					  description
					  args {
						name
						description
					  }
					}
					fields {
					  name
					}
				  }
				}
			  }
			}
		  }
		}
	  }`
