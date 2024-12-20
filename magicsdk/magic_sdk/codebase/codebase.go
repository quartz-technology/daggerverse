package codebase

import (
	"fmt"
	"os"
)

type Codebase struct {
	Path string

	dir []os.DirEntry
}

func New(path string) (*Codebase, error) {
	fmt.Printf("Reading codebase at %s\n", path)

	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	return &Codebase{
		Path: path,
		dir:  dir,
	}, nil
}

func (c *Codebase) LookupFile(name string) (*os.File, error, bool) {
	for _, entry := range c.dir {
		fmt.Printf("Checking file %s\n", entry.Name())

		if entry.Name() == name {
			file, err := os.Open(c.Path + "/" + entry.Name())
			if err != nil {
				return nil, err, false
			}

			return file, nil, true
		}
	}

	return nil, nil, false
}