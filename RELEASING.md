# Daggerverse release process

## Requirements

- [changie](https://changie.dev): a changelog management tool.
- [GitHub CLI](https://cli.github.com/): GitHub utility tools.

## Verify your `main` branch is up to date

```shell
git checkout main

git pull

git fetch --all --tags
```

## Store the current version

```shell
export VERSION=$(changie latest)

# echo $VERSION
```

## Update CHANGELOG

```shell
git checkout -b release/bump-$VERSION

changie batch patch
changie merge

# Open the Changelog update PR and merge it if everything is okay.
```

💡 Replace `batch` by minor or major if it's another version.

## Tag the release

```shell
git tag origin $VERSION
git push origin $VERSION

gh release create "Quartz Daggerverse - $VERSION" --draft --verify-tag --title "Quartz Daggerverse - $VERSION" --notes-files .changes/$VERSION.md
``

This will trigger the workflow to publish all modules to daggerverse.

