# Bun server

## Getting started

You don't need any tools locally except [`Dagger`](https://dagger.io), every operation such as infrastructure providing, testing, running the application etc.. will be performed using `Dagger`.

To explore the different commands, please run:

```shell
$ dagger functions                      
Name    Description
app     Provide the source code and a Bun container with dependencies installed
build   Return the built Bun application
ci      Execute build, lint, and test in concurrency
infra   Access specific infrastructure commands
lint    Execute lint on the app
run     Start the application on port 8080 with the development stack
stack   Create the development stack for the Bun application
test    Runs the integration tests of the application in an isolate environment
```

## Demos

This demo available in [`demos.mov`](./demos.mov) shows how to run your CI and your local application at the same time.
This spins up 2 isolated development environments and runs concurrently inside Dagger.

We then execute Postman queries to verify the dev application while the CI is running on its own.

![Watch the video](https://imgur.com/jFZIFy7)

## Run the project

You can start the development infrastructure with the application and its services (PostgreSQL, Redis) using:

```shell
dagger call app --source=. run up
```

This is going to call the `app` and `run` functions defined in [`dagger/src/index.ts`](./dagger/src/index.ts)

## Run the CI

You can run the integration tests, build and lint in concurrency using:

```shell
dagger call app --source=. ci
```