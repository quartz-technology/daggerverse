package main

import "fmt"

// Database returns a ready to run Postgres container with all configuration applied.
func (p *Postgres) Database() (*Container, error) {
	startOpts := []string{}

	ctr := dag.
		Container().
		From(fmt.Sprintf("postgres:%s", p.Version))

	if p.Cache {
		ctr = ctr.WithMountedCache("/var/lib/postgresql/data", dag.CacheVolume("pg-data"))
	}

	ctr = ctr.
		WithSecretVariable("POSTGRES_USER", p.User).
		WithSecretVariable("POSTGRES_PASSWORD", p.Password)

	// Set database name
	if p.Name != "" {
		ctr = ctr.WithEnvVariable("POSTGRES_DATABASE", p.Name)
	}

	// Set config files
	if p.ConfigFile != nil {
		ctr = ctr.WithFile("/etc/postgresql/postgresql.conf", p.ConfigFile)
		startOpts = append(startOpts, "-c", "config_file=/etc/postgresql/postgresql.conf")
	}

	// Set init scripts
	if p.InitScripts != nil {
		ctr = ctr.WithMountedDirectory("/docker-entrypoint-initdb.d", p.InitScripts)
	}

	// Apply start opts
	ctr = ctr.
		WithExposedPort(p.Port).
		WithExec(startOpts)

	return ctr, nil
}
