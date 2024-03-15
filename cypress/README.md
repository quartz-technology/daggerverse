# Cypress module

A Cypress module to run your Cypress test on the project source.

## Usage

:bulb: You can provide a Node project with a `build` command that outputs the source on dist and the module will figure out a way to build and launch it.

### Run Cypress test on your project

```shell
dagger -m github.com/quartz-technology/daggerverse/cypress --source=<path> call run
```

Made with ❤️ by Quartz.