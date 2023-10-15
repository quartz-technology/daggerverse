package main

import (
	"context"
	"fmt"
	"strconv"
)

type Redis struct {
	port     int
	version  string
	password *Secret
}

func (r *Redis) WithPort(port int) *Redis {
	r.port = port

	return r
}

func (r *Redis) WithVersion(version string) *Redis {
	r.version = version

	return r
}

func (r *Redis) WithPassword(password *Secret) *Redis {
	r.password = password

	return r
}

func (r *Redis) Server(ctx context.Context) (*Container, error) {
	ctr := dag.
		Container().
		From(fmt.Sprintf("bitnami/redis:%s", r.version)).
		WithMountedCache("/bitnami/redis/data", dag.CacheVolume("redis-data"))

	if r.password != nil {
		password, err := r.password.Plaintext(ctx)
		if err != nil {
			return nil, err
		}

		ctr = ctr.WithEnvVariable("REDIS_PASSWORD", password)
	} else {
		ctr = ctr.WithEnvVariable("ALLOW_EMPTY_PASSWORD", "yes")
	}

	if r.port != 0 {
		r.port = 6379
	}

	ctr = ctr.
		WithEnvVariable("REDIS_PORT_NUMBER", strconv.Itoa(r.port)).
		WithExposedPort(r.port)

	return ctr, nil
}

func (r *Redis) Cli(ctx context.Context) (*Container, error) {
	ctr, err := r.Server(ctx)
	if err != nil {
		return nil, err
	}

	ctr = ctr.WithEntrypoint([]string{"redis-cli"})

	return ctr, nil
}
