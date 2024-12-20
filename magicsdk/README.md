# Magic SDK

## Usage

:warning: This is experimental, for now it only supports `Go` project and `Dockerfile`.

```shell
# Go to your project root

# Init magicSDK
dagger init --sdk github.com/quartz-technology/daggerverse/magicsdk --name=test .

# Explore functions that are available
dagger functions
```

## Current integrations

### Docker

`dagger call docker build`: use your project Dockerfile to build a container and return it.

### Go

`dagger call go container`: creates a development environment for your project inside a container and return it.