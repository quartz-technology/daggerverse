# Launcher

Publish a dagger module using a dagger module.

### Usage

```go
dag.Launcher().Publish(ctx, repo(), LauncherPublishOpts{Path: "<path to the module>"})
```

```shell
dagger -m github.com/quartz-technology/daggerverse/launcher call --module <repository> --path <path to the module> 
```

:warning: You must first push your module to git to let daggerverse index it.

Made with ❤️ by Quartz.
