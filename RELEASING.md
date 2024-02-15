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

## Update CHANGELOG

```shell
# Batch changelor prior to opening the PR so the version
# is bump in the right order
changie batch patch
changie merge

export VERSION=$(changie latest)
git checkout -b release/bump-$VERSION

# Open the Changelog update PR and merge it if everything is okay.
```

💡 Replace `batch` by minor or major if it's another version.

## Tag the release

```shell
# Make sure you merged the changelog PR first
git checkout main
git pull

export VERSION=$(changie latest)


git tag $VERSION origin 
git push origin $VERSION

gh release create "$VERSION" --draft --verify-tag --title "$VERSION" --notes-file .changes/$VERSION.md
```

This will trigger the workflow to publish all modules to daggerverse.

