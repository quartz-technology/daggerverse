package main

import "fmt"

type Postgres struct {
	User        *Secret
	Password    *Secret
	Database    string
	Port        int
	Version     string
	ConfigFile  *File
	InitScripts map[string]*File
}

// WithCredential adds a user and a password configuration to the postgresSQL
// database.
// The values will be set in container as secret variables: `POSTGRES_USER`
// `POSTGRES_PASSWORD`.
func (p *Postgres) WithCredential(user *Secret, password *Secret) *Postgres {
	p.User = user
	p.Password = password

	return p
}

// WithPort configs an exposed port on the container.
func (p *Postgres) WithPort(port int) *Postgres {
	p.Port = port

	return p
}

// WithDatabase sets the name of the database that will be created on start.
// It will be set in the container as `POSTGRES_DB`, if it's not set, the
// database's name will be the user's one.
func (p *Postgres) WithDatabase(name string) *Postgres {
	p.Database = name

	return p
}

// WithVersion sets the version of postgresql to pull from the registry.
func (p *Postgres) WithVersion(version string) *Postgres {
	p.Version = version

	return p
}

// WithConfigFile adds an extra config file to the postgres database.
// This file will be copied in the container to
// `/usr/share/postgresql/postgresql.conf`
func (p *Postgres) WithConfigFile(file *File) *Postgres {
	p.ConfigFile = file

	return p
}

// WithInitScript adds a script to execute when the database is started.
// You can call this function multiple times to add multiple scripts.
// These scripts are stored in a map, so it's recommended to name with a numeric
// value at the beginning to make sure they are executed in the correct order.
// For example `1-init.sql`, `2-new-tabs.sql`...
//
// Theses files will be copied to `/docker-entrypoint-initdb.db`
func (p *Postgres) WithInitScript(name string, script *File) *Postgres {
	p.InitScripts[name] = script

	return p
}

// Start returns a ready to run Postgres container with all configuration applied.
func (p *Postgres) Start() (*Container, error) {
	startOpts := []string{}

	ctr := dag.
		Container().
		From(fmt.Sprintf("postgres:%s", p.Version)).
		WithMountedCache("/var/lib/postgresql/data", dag.CacheVolume("pg-data"))

	// Set credential
	if p.User == nil || p.Password == nil {
		return nil, fmt.Errorf("start error: username and password required, call Withcredential before Start")
	}

	ctr = ctr.
		WithSecretVariable("POSTGRES_USER", p.User).
		WithSecretVariable("POSTGRES_PASSWORD", p.Password)

	// Set database name
	if p.Database != "" {
		ctr = ctr.WithEnvVariable("POSTGRES_DATABASE", p.Database)
	}

	if p.ConfigFile != nil {
		ctr = ctr.WithFile("/etc/postgresql/postgresql.conf", p.ConfigFile)
		startOpts = append(startOpts, "-c", "config_file=/etc/postgresql/postgresql.conf")
	}

	// Set init scripts
	for filename, file := range p.InitScripts {
		ctr = ctr.WithMountedFile(
			fmt.Sprintf("/docker-entrypoint-initdb.d/%s", filename),
			file,
		)
	}

	// Apply start opts
	ctr = ctr.
		WithExposedPort(p.Port).
		WithExec(startOpts)

	return ctr, nil
}
