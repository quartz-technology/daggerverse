package main

import "fmt"

// Server returns a Container with a Minio server ready to be started.
// If set, the server adds authentication with username/password
// (access/secret keys).
//
// By default, the server listens on port 9000 and console on port
// 9001, but it can be sets to another value with methods.
func (m *Minio) Server() *Container {
	ctr := dag.
		Container().
		From(fmt.Sprintf("quay.io/minio/minio:%s", m.Version))

	// Setup port with given or default values
	if m.ServerPort == 0 {
		m.ServerPort = 9000
	}

	if m.ConsolePort == 0 {
		m.ConsolePort = 9001
	}

	ctr = ctr.
		WithExposedPort(m.ServerPort).
		WithExposedPort(m.ConsolePort)

	// Setup credential
	if m.Username != nil && m.Password != nil {
		ctr = ctr.
			WithSecretVariable("MINIO_ROOT_USER", m.Username).
			WithSecretVariable("MINIO_ROOT_PASSWORD", m.Password)
	}

	// Set cache if set to true
	if m.Cache {
		ctr = ctr.WithMountedCache("/data", dag.CacheVolume("minio-data"))
	}

	return ctr
}
