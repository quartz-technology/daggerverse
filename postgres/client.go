package main

import (
	"fmt"
	"strconv"
)

// Client returns a configured client with user and password set.
// Note: if you want to set the host, you can append this option on usage
func (p *Postgres) Client() (*Container, error) {
	cliOpts := []string{"/bin/sh", "-c", "psql"}

	ctr := dag.
		Container().
		From(fmt.Sprintf("postgres:%s", p.Version))

	// Set credential
	if p.User == nil || p.Password == nil {
		return nil, fmt.Errorf("start error: username and password required, call Withcredential before Start")
	}

	ctr = ctr.
		WithSecretVariable("PGPASSWORD", p.Password).
		WithSecretVariable("PGUSER", p.User)

	cliOpts = append(cliOpts, "-U", "$PGUSER")

	if p.Name != "" {
		cliOpts = append(cliOpts, "-d", p.Name, "-p", strconv.Itoa(p.Port))
	}

	ctr = ctr.
		WithEntrypoint(cliOpts)

	return ctr, nil
}
