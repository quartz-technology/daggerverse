# Docker SDK

A magic SDK that automatically provides a custom function to build your project's Dockerfile.

Any argument, secret, or stages defined in your Dockerfile will be automatically exposed as parameter to your function.

## Usage 

:warning: This is experimental, please open an issue if you encounter any issue.

```shell
# Go to your project root (or any directory that contains a Dockerfile)

# Init Dagger module with the Docker SDK 
# The SDK will look in the `source` directory for a Dockerfile, so keep it to `.`
dagger init --sdk github.com/quartz-technology/daggerverse/docker_sdk --name=test --source=.

# Build a container of project's application using your dockerfile and start a terminal inside it
dagger call docker build terminal 
```

## Functions

### Build

Build a container of your project's application using your dockerfile and return it.

**Supported arguments**
- `platform`: The platform to build the container for (default to your host platform).
- `dockerfile`: The path to the Dockerfile to use (default to `Dockerfile` or the first Dockerfile found in the project).
- `target`: The target stage to build (optional and will build the last stage by default).
- `buildArgs`: A list of build arguments to pass to the build (optional).
- `secrets`: A list of secrets to pass to the build (optional).

## Example

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