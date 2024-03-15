package main

import (
	"fmt"
	"strconv"
)

// Server returns a new container running Redis a redis Server.
func (r *Redis) Server() (*Container, error) {
	ctr := dag.
		Container().
		From(fmt.Sprintf("bitnami/redis:%s", r.Version)).
		WithUser("root")

	if r.Cache  {
		ctr = ctr.WithMountedCache("/bitnami/redis/data", dag.CacheVolume("redis-data"))
	}

	if r.Password != nil {
		ctr = ctr.WithSecretVariable("REDIS_PASSWORD", r.Password)
	} else {
		ctr = ctr.WithEnvVariable("ALLOW_EMPTY_PASSWORD", "yes")
	}

	ctr = ctr.
		WithEnvVariable("REDIS_PORT_NUMBER", strconv.Itoa(r.Port)).
		WithExposedPort(r.Port)

	return ctr, nil
}
