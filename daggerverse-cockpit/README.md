# Daggerverse cockpit

A module to help you manage your Daggerverse repository

## Usage

### Generate a usage boilerplate

```shell
dagger -m github.com/quartz-technology/daggerverse/daggerverse-cockpit call usage_generator --module=.
```

### Publish all your modules

```shell
dagger -m github.com/quartz-technology/daggerverse/daggerverse-cockpit call publish --repository=. --exclude=<exclude paths> 
```

💡 Use `--dryRun` option if you want to check if the module selects the right modules to publish.

💡 It's possible to use the CLI to use this module but you can benefit from its full power if you integrate it on a CI.

You can use the following GitHub Actions to automatically publish your modules on releases

```yaml
name: Release Daggerverse modules

on:
  release:
    types: [published]

jobs:
  release:
    name: Release Daggerverse
    runs-on: ubuntu-latest
    env:
      DAGGER_CLOUD_TOKEN: ${{ secrets.DAGGER_CLOUD_TOKEN }}

    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Install dagger
        run: |
          cd /usr/local
          curl -L https://dl.dagger.io/dagger/install.sh | sh
          dagger version

      - name: Run publishing
        run: |
          dagger -m github.com/quartz-technology/daggerverse/daggerverse-cockpit call publish --repository=. --exclude=<dir to exclude>  
```

Made with ❤️ by [Quartz](https://quartz.technology).