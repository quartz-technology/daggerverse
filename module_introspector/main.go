package main

import (
	"context"
	"dagger/module-introspector/internal/dagger"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

type ModuleIntrospector struct{}

// Returns a container that echoes whatever string argument is provided
func (m *ModuleIntrospector) Introspect(ctx context.Context, moduleSource *dagger.Directory) (*dagger.File, error) {
	err := moduleSource.AsModule().Initialize().Serve(ctx)
	if err != nil {
		return nil, err
	}

	var resp response
	err = dag.GraphQLClient().MakeRequest(ctx, &graphql.Request{
		Query:  query,
		OpName: "IntrospectionQuery",
	}, &graphql.Response{
		Data: &resp,
	})
	if err != nil {
		return nil, fmt.Errorf("introspection query: %w", err)
	}

	jsonInstrospection, err := resp.AsJSON()
	if err != nil {
		return nil, fmt.Errorf("introspection query: %w", err)
	}

	return dag.
		Directory().
		WithNewFile("introspection.json", jsonInstrospection).
		File("introspection.json"), nil
}