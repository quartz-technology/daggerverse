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