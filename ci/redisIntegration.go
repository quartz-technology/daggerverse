package main

import (
	"context"
	"fmt"
)

func (i *IntegrationTest) Redis(ctx context.Context) error {
	password := dag.SetSecret("redis-password", "foo123")

	redisCtr := dag.Redis().
		WithVersion("latest").
		WithPassword(password)

	server := redisCtr.Server()

	cli := redisCtr.Cli(server.AsService())

	_, err := cli.Set("foo", "bar").Sync(ctx)
	if err != nil {
		return err
	}

	value, err := cli.Get(ctx, "foo")
	if err != nil {
		return err
	}

	if value != "bar\n" {
		return fmt.Errorf("value are not matching, expecting: '%s', got: '%s'", "bar", value)
	}

	return nil
}
