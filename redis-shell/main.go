package main

type RedisShell struct{}

func (m *RedisShell) Cli() *Container {
	password := dag.SetSecret("redis-password", "foo123")

	redisCtr := dag.Redis().
		WithVersion("latest").
		WithPassword(password)

	server := redisCtr.Server()

	cli := redisCtr.Cli(server.AsService())

	return cli.Container()
}
