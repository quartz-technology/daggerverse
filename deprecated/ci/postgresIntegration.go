package main

import (
	"context"
	"time"
)

func (i *IntegrationTest) Postgres(ctx context.Context) error {
	password := dag.SetSecret("postgres_password", "XxXSecretPwd")
	user := dag.SetSecret("postgres_user", "daggerverse-test")

	postgres := dag.
		Postgres().
		WithVersion("13").
		WithCredential(user, password).
		WithDatabaseName("daggerverse").
		WithCache(false)

	db := postgres.Database()

	example := dag.
		Git("github.com/prisma/prisma-examples").
		Branch("latest").
		Tree().
		Directory("databases/postgresql-supabase")

	nodeCtr := dag.Node().
		WithVersion("20-alpine3.17").
		WithSource(example).
		WithNpm().
		Install([]string{}).
		Container().
		WithServiceBinding("postgres", db.AsService())

	dbUrl := dag.SetSecret("postgres_url", "postgres://daggerverse-test:XxXSecretPwd@postgres:5432/daggerverse")

	// Migrate and exec dev to send a query and verify the connection
	_, err := nodeCtr.
		WithSecretVariable("DATABASE_URL", dbUrl).
		WithEnvVariable("CACHE_BURST", time.DateTime).
		WithExec([]string{"npx", "prisma", "migrate", "dev", "--name", "init"}, ContainerWithExecOpts{SkipEntrypoint: true}).
		WithExec([]string{"run", "dev"}).Stdout(ctx)
	if err != nil {
		return nil
	}

	return nil
}
