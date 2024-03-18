package main

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/iancoleman/strcase"
)

var usageFuncMap = template.FuncMap{
	"BackTick": func() string {
		return "```"
	},
	"ToPascalCase": strcase.ToCamel,
	"ToSnakeCase":  strcase.ToSnake,
}

// UsageGenerator generates a simple usage documentation for a module.
//
// This function is still in developement, it's a simple utility functions
// to simplify modules maintaing.
//
// Use it as a starting point for your module README but please, do not
// consider it as a README generator.
//
// Example usage: dagger call usage-generator --module=. -o USAGE.md
func (d *DaggerverseCockpit) UsageGenerator(
	ctx context.Context,

	module *Directory,
) (*File, error) {
	introspection, err := d.introspectModule(ctx, module)
	if err != nil {
		return nil, fmt.Errorf("failed to introspect the module: %w", err)
	}

	tmpl, err := template.New("usage").Funcs(usageFuncMap).Parse(usageTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to create the usage template: %w", err)
	}

	var result bytes.Buffer
	if err = tmpl.Execute(&result, introspection[0]); err != nil {
		return nil, fmt.Errorf("failed to execute the usage template: %w", err)
	}

	return dag.Directory().WithNewFile("USAGE.md", result.String()).File("USAGE.md"), nil
}
