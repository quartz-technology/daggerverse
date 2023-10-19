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

	cli := redisCtr.
		Cli().
		WithServiceBinding("redis", server.AsService()).
		WithEntrypoint([]string{"redis-cli", "-h", "redis"})

	// Set key
	_, err := cli.WithExec([]string{"SET", "foo", "bar"}).Sync(ctx)
	if err != nil {
		return err
	}

	// Get key
	value, err := cli.WithExec([]string{"GET", "foo"}).Stdout(ctx)
	if err != nil {
		return err
	}

	if value != "bar" {
		return fmt.Errorf("value are not matching, expecting: %s, got: %s", "foo", value)
	}

	return nil
}
