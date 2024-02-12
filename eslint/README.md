# Eslint module

A Typescript linter module for standalone typescript files and Node projects.

## Usage

:bulb: If the given files have no `package.json`, `eslintrc` config, or `tsconfig.json`, a default
one will be plugged into the container.

### Open a shell inside the linter container

```shell
dagger -m github.com/quartz-technology/daggerverse/eslint --files=<path> call container terminal
```

:bulb: It's a practical way to test the module / debug any issue that may happen inside the container.

### Run the linter on Typescript files

```shell
dagger -m github.com/quartz-technology/daggerverse/eslint --files=<path> call run stdout
```

Made with ❤️ by Quartz.