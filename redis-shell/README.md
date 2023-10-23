# Redis shell

⚠️ This module is an experimentation of `dagger shell` command, it's not made for any production usage.

This module is importing [`redis`](../redis) module to start a shell in order to interact with the
actual server.

### Requirement

- Dagger [v0.9.0](https://github.com/dagger/dagger/releases/tag/v0.9.0) CLI and engine installed on your host

### Getting Started

```sh
dagger shell cli --entrypoint /bin/sh
```

Made with ❤️ by Quartz.