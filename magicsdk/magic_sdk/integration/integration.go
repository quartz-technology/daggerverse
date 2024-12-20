package integration

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"dagger.io/magicsdk/codebase"
	"dagger.io/magicsdk/invocation"
)

type Integrations map[string]Integration

type Integration interface {
	Exist() bool

	New(dir *dagger.Directory) Integration

	TypeDef() *dagger.TypeDef

	Invoke(ctx context.Context, invocation *invocation.Invocation) (_ any, err error)
}

type integrationFunc func(code *codebase.Codebase) (Integration, error)

var integrationsFuncs = map[string]integrationFunc{
	"Docker": DockerIntegration,
}

func LoadIntegrations(code *codebase.Codebase) (Integrations, error) {
	integrations := make(map[string]Integration)
	
	for name, integrationFct := range integrationsFuncs {
		fmt.Println("trying to load integration", name)

		integration, err := integrationFct(code)
		if err != nil {
			return nil, fmt.Errorf("failed to load integration: %w", err)
		}

		if integration.Exist() {
			fmt.Println("integration", name, "exists in that project")

			integrations[name] = integration
		} else {
			fmt.Println("integration", name, "doesn't exist in that project")
		}
	}

	return integrations, nil
}
