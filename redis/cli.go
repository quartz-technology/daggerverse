package main

import (
	"context"
	"redis/internal/dagger"
)

type Cli struct {
	Ctr *dagger.Container
}

// CLI returns a new container running the Redis CLI connected to a redis Service.
func (r *Redis) Cli(
	// The Redis server to connect to.
	server *dagger.Service,
) (*Cli, error) {
	ctr, err := r.Server()
	if err != nil {
		return nil, err
	}

	entrypointCmd := []string{"redis-cli", "-h", "redis"}
	if r.Password != nil {
		ctr = ctr.WithSecretVariable("REDISCLI_AUTH", r.Password)
	}

	ctr = ctr.
		WithServiceBinding("redis", server).
		WithEntrypoint(entrypointCmd)

	return &Cli{
		Ctr: ctr,
	}, nil
}

func (c *Cli) Container() *dagger.Container {
	return c.Ctr
}

func (c *Cli) Set(key, value string) *dagger.Container {
	return c.Ctr.WithExec([]string{"SET", key, value}, dagger.ContainerWithExecOpts{UseEntrypoint: true})
}

func (c *Cli) Get(ctx context.Context, key string) (string, error) {
	return c.Ctr.WithExec([]string{"GET", key}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).Stdout(ctx)
}
