package main

import (
	"fmt"
	"strconv"
)

func (r *Redis) Server() (*Container, error) {
	ctr := dag.
		Container().
		From(fmt.Sprintf("bitnami/redis:%s", r.Version)).
		WithUser("root").
		WithMountedCache("/bitnami/redis/data", dag.CacheVolume("redis-data"))

	if r.Password != nil {
		ctr = ctr.WithSecretVariable("REDIS_PASSWORD", r.Password)
	} else {
		ctr = ctr.WithEnvVariable("ALLOW_EMPTY_PASSWORD", "yes")
	}

	if r.Port == 0 {
		r.Port = 6379
	}

	ctr = ctr.
		WithEnvVariable("REDIS_PORT_NUMBER", strconv.Itoa(r.Port)).
		WithExposedPort(r.Port)

	return ctr, nil
}
