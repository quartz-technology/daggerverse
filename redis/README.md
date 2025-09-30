# Redis

A simple module to start and interact with a Redis service.

| Command                | Done |
|------------------------|------|
| Setup a Redis server   | ✅    |
| Setup a Redis CLI      | ✅    |
| Set a key              | ✅    |
| Get a key              | ✅    |
| Configure Redis server | ✅    |
| Setup authentication   | ✅    |
| Clusters               | ❌    |

## Usage

### Create a Redis server

```shell
dagger -m github.com/quartz-technology/daggerverse/redis call server up
```

#### Test the server from local
```shell
dagger -m github.com/quartz-technology/daggerverse/redis call cli --server=tcp://localhost:6379 set --key=foo --value=bar
dagger -m github.com/quartz-technology/daggerverse/redis call cli --server=tcp://localhost:6379 get --key=foo
```

### Create a Redis client

⚠️ Only available as code since you need to provide a Redis server to the CLI.

Made with ❤️ by Quartz.
