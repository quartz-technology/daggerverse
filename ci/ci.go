package main

import (
	"os"
	"path/filepath"
)

type Ci struct{}

func repo() *Directory {
	return dag.
		Host().
		Directory(root())
}

func source() *Directory {
	return dag.
		Host().
		Directory(root(), HostDirectoryOpts{
			Include: []string{"**/*.go", "**/go.mod", "**/go.sum", ".golangci.yml"},
		})
}

func daggerRepository() *Directory {
	return dag.
		Git("github.com/dagger/dagger").
		Branch("main").
		Tree()
}

// TODO: fix .. restriction
func root() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(wd, "..")
}
