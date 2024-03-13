# Postgres

A simple module to start a PostgreSQL service.

| Command                                   | Done |
|-------------------------------------------|------|
| Setup a Postgres database                 | ✅    |
| Setup credential                          | ✅    |
| Configure Postgres database               | ✅    |
| Connect to Postgres database from service | ✅    |
| Enable or disable cache ddata             | ✅    |

:warning: PSQL cannot be used for now because it requires a Socket connection and it's currently
not possible to bind a socket from a container.

## Usage

### Create a client to a Postgres database

```shell
dagger call --user john --password 123456 --name test client 
```

### Create a Postgres database

```shell
dagger call --user john --password 123456 --name database 
```


Made with ❤️ by Quartz.
