package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

type DaggerverseCockpit struct{}

// Publish loop through all your directory that contains a `dagger.json`
// and publish them to the daggerverse.
func (d *DaggerverseCockpit) Publish(
	ctx context.Context,

	// The repository that contains your dagger modules
	repository *Directory,

	// Excluse some directories from publishing
	// It's useful if you use this module from Dagger CLI.
	// +optional
	exclude []string,

	// Only returns the path of the modules that shall be published
	// +optional
	dryRun bool,
) ([]string, error) {
	// Loop through all the directories and check for `dagger.json` files
	daggerJSONPaths, err := repository.Glob(ctx, "**/dagger.json")
	if err != nil {
		return nil, fmt.Errorf("could not retrieve your dagger modules: %w", err)
	}

	// Exclude modules that should not be published 
	paths := excludeModules(daggerJSONPaths, exclude)

	// Remove `dagger.json` from path
	for i, path := range paths {
		paths[i] = filepath.Dir(path)
	}

	if dryRun {
		return paths, nil
	}

	// Publish the modules to the daggerverse
	eg, ctx := errgroup.WithContext(ctx)

	for i, path := range paths {
		i := i
		path := path

		eg.Go(func() error {
			result, err := dag.DaggerPublisher().Publish(ctx, repository, DaggerPublisherPublishOpts{Path: path})
			if err != nil {
				return err
			}

			// Retrieve the url of the module
			logs := strings.Split(result, "\n")
			url := strings.TrimSpace(logs[len(logs)-2])

			// Add it to the output path
			paths[i] += fmt.Sprintf(" published to: %s", url)

			return nil
		})
	}

	if err = eg.Wait(); err != nil {
		return nil, fmt.Errorf("could not publish your dagger modules: %w", err)
	}

	return paths, nil
}
