package main

import "context"

type Cli struct {
	Ctr *Container
}

func (r *Redis) Cli(server *Service) (*Cli, error) {
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

func (c *Cli) Container() *Container {
	return c.Ctr
}

func (c *Cli) Set(key, value string) *Container {
	return c.Ctr.WithExec([]string{"SET", key, value})
}

func (c *Cli) Get(ctx context.Context, key string) (string, error) {
	return c.Ctr.WithExec([]string{"GET", key}).Stdout(ctx)
}
