# Docker SDK

A SDK that automatically provides a custom function to build your project's Dockerfile or run your docker-compose services.

## Dockerfile

Any argument, secret, or stages defined in your Dockerfile will be automatically exposed as parameter to your function.

### Usage 

:warning: This is experimental, please open an issue if you encounter any issue.

```shell
# Go to your project root (or any directory that contains a Dockerfile)

# Init Dagger module with the Docker SDK 
# The SDK will look in the `source` directory for a Dockerfile, so keep it to `.`
dagger init --sdk github.com/quartz-technology/daggerverse/docker_sdk --name=test --source=.

# Build a container of project's application using your dockerfile and start a terminal inside it
dagger call docker build terminal

# Start all services defined in the docker-compose file
dagger call docker compose all up

# Start a single service by its name
dagger call docker compose redis up
```

### Functions

#### Build

Build a container of your project's application using your dockerfile and return it.

**Supported arguments**
- `platform`: The platform to build the container for (default to your host platform).
- `dockerfile`: The path to the Dockerfile to use (default to `Dockerfile` or the first Dockerfile found in the project).
- `target`: The target stage to build (optional and will build the last stage by default).
- `buildArgs`: A list of build arguments to pass to the build (optional).
- `secrets`: A list of secrets to pass to the build (optional).

### Example

```Dockerfile
ARG BASE_IMAGE=golang:1.23.2-alpine

FROM ${BASE_IMAGE} AS app

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/main .

FROM golang:1.23.2-alpine AS runtime

ARG BIN_NAME

WORKDIR /runtime

COPY --from=app /app/main /runtime/${BIN_NAME}

RUN --mount=type=secret,id=my-super-secret \
  cat /run/secrets/my-super-secret > /runtime/secret.txt

ENTRYPOINT ["/runtime/${BIN_NAME}"]
```

If you display the helper of the `build` function, you see that build arguments and secrets declared
in the dockerfile are automatically exposed as arguments to the function (with default value if it exists).
For stages, it's exposed as an enumeration to enforce validation.

```shell
dagger call docker build --help

ARGUMENTS
      --bin-name string          Set BIN_NAME build argument [required]
      --my-super-secret Secret   Set my-super-secret secret [required]
      --base-image string        Set BASE_IMAGE build argument (default "golang:1.23.2-alpine")
      --dockerfile string        Path to the Dockerfile to use. (default "Dockerfile")
      --platform Platform        Platform to build. (default linux/arm64)
      --target app,runtime       Target stage to build.
```

Example of command:

```shell
dagger call docker build --bin-name foooo --my-super-secret env:PATH terminal

$ ls
foooo       secret.txt
```

## Docker Compose

If a docker compose file (`docker-compose.[yaml|yml]`, `compose.[yaml|yml]`) is present in the current directory, it will be parsed and accessible
with `dagger call docker compose`.

Any services defined in your docker compose file will be registered as callable functions (by its name) and returns the service's container.
All properties of a service will be used to build the service's container and may be settable as arguments to the callable function.

### Supported properties

| Property      | Description                                           | Settable from Dagger CLI                      | 
|---------------|-------------------------------------------------------|-----------------------------------------------|
| `image`       | The image to use for the service                      | Yes                                           |        
| `build`       | The build context and options for the service         | No                                            |
| `workdir`     | The workdir to set for this container                 | No                                            | 
| `command`     | Command to run inside the service container           | Yes                                           |
| `entrypoint`  | Override the default entrypoint of the container      | Yes                                           |
| `environment` | Environment variables to set in the container         | Yes ([details here](#environment-variables))  |
| `ports`       | Ports to expose from the container                    | Yes                                           |
| `volumes`     | Volumes to mount in the container                     | Yes ([details here](#volumes))                |
| `depends_on`  | Services to depend on for the service to start        | No  ([details here](#depends-on))             | 

#### Environment variables

Environment variables are defined as a map of key/value pairs.

For example:

```yaml
my-service:
  environment:
    -FOO=bar
    -BAZ
```

Can be set as arguments to the callable function:

```shell
dagger call docker compose my-service --foo bar
```

If a value if already defined in the compose file, it will be register as the variable's default value.

If no value is defined, it will be registered as an optional secret value.

A secret value can be set as an argument to the callable function:

```shell
SECRET=foo dagger call docker compose my-service --baz env:SECRET
```

#### Volumes

Volumes are defined as a list of strings.

For example:

```yaml
my-service:
  volumes:
    - data:/app/data  # Cache volume
    - .:/code         # Mountable volume
    - ./my-file.txt   # Mountable file
```

If a volume is defined as a cache volume, it will be mount as cache volume to the container.

A mountable volume will be mounted as a directory or file depending on the type of the volume and settable as an argument to the callable function:

```shell
dagger call docker compose my-service --current-directory .
```

By default, the volume argument will be set to the `defaultPath` as defined in the volume declaration. (e.g., `.:/code` will mount the current directory to `/code` by default).

If the mount path is set to `.`, it will be aliased by `current-directory`.

#### Depends on

The `depends-on` argument is a list of other services that must be bind to the called service.

```yaml
services:
  my-service:
    depends-on:
      - my-other-service
      - my-other-service2
```

If dependent services also depends on other services, they will be also be bound to the dependent service.

Any argument of these services will be available as argument in the CLI, prefixed by the service name.

For example, if `my-other-service` has an argument `my-arg`, it will be available as `--my-other-service-my-arg` in the CLI.

### Example

```yaml
services:
  backend:
    build: ./backend
    environment:
      MESSAGE: "Hello from Docker Compose"
    ports:
      - "8080:8080"

  gateway:
    build: ./gateway
    environment:
      - MESSAGE="test"
    ports:
      - "8081:8081"
    depends_on:
      - backend

  redis:
    image: bitnami/redis:latest
    environment:
      - REDIS_PASSWORD
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/bitnami/redis/data

volumes:
  redis_data:
```

#### All

Start all services with the function `all`:

```shell
$ dagger call docker compose all --help

ARGUMENTS
      --backend-message string        Set environment variable MESSAGE (default "Hello from Docker Compose")
      --gateway-message string        Set environment variable MESSAGE (default "test")
      --redis-image string            Image to use for the service (default "bitnami/redis:latest")
      --redis-redis-password Secret   Set secret environment variable REDIS_PASSWORD

$ PASSWORD=test dagger call docker compose all --redis-redis-password env:PASSWORD up

# Services will be accessible at http://localhost:8080, http://localhost:8081, and http://localhost:6379
```

:bulb: All functions will have their arguments prefixed by their service name.

#### Service

Start a single service by its name:

```shell
dagger call docker compose gateway --help

ARGUMENTS
      --backend-message string   Set environment variable MESSAGE (default "Hello from Docker Compose")
      --message string           Set environment variable MESSAGE (default "test")

dagger call docker compose gateway up

# Service will be accessible at http://localhost:8081 (Only gateway service is exposed to host in that case)
```