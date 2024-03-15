# Nginx module

A Nginx module to expose HTML files through a proxy server.

## Usage

:bulb: You can either provide your Nginx configuration or use the default one.

:warning: This module only targets exposing HTML Website, it's not done for load balancing nor
complex configuration.

### Expose your built source on a server

```shell
dagger -m github.com/quartz-technology/daggerverse/nginx --files=<path> call expose up
```

Made with ❤️ by Quartz.