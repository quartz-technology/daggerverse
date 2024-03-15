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
dagger -m github.com/quartz-technology/daggerverse/redis call server expose up
```

### Create a Redis client

⚠️ Only available as code since you need to provide a Redis server to the CLI.

Made with ❤️ by Quartz.
