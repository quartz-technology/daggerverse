# Module Introspector

This module is inspired from [Helder's codegen module](https://daggerverse.dev/mod/github.com/helderco/daggerverse/codegen).

It is used to introspect the GraphQL schema of a given module and returns its introspection query response as a JSON file.

## Usage

### Introspect a remote module

```shell
dagger call introspect --module-source https://github.com/kpenfound/dagger-modules.git\#main:proxy -o introspection-result.json
``

### Introspect a local module

```shell
dagger call introspect --module-source ./my-module -o introspection-result.json
```

This can then be used with https://github.com/dagger/dagger/pull/9485 to generate client bindings.