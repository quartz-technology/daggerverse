# Dagger CLI module

A module to use the Dagger CLI

💡 This module is useful to integrate Dagger inside more complex CI.

## Usage

### Publish a module 

```shell
dagger call -m github.com/quartz-technology/daggerverse/dagger-cli --module=. call publish
```

### Call a module

```shell
dagger call -m github.com/quartz-technology/daggerverse/dagger-cli --module=<path> call --args=<args>
```

Made with ❤️ by Quartz.