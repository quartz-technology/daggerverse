package finder

import (
	"fmt"
	"os"
	"path/filepath"
)

// Finder helps locate files.
//
// It operates within a specified root directory and its entries.
type Finder struct {
	// dirPath specifies the root directory for searches and checks.
	dirPath string

	// dir holds the directory entries within dirPath.
	dir []os.DirEntry
}

// New creates a Finder with a specified directory path and entries.
func New(dirPath string, dir []os.DirEntry) *Finder {
	return &Finder{
		dirPath: dirPath,
		dir:     dir,
	}
}

// FindFileFromPattern searches for files matching any pattern.
//
// Returns the full path of the first match and true, or an empty string
// and false if no matches are found.
func (f *Finder) FindFileFromPattern(patterns []string) (string, bool) {
	for _, entry := range f.dir {
		for _, pattern := range patterns {
			matches, err := filepath.Match(pattern, entry.Name())
			if err != nil {
				fmt.Printf("failed to match pattern %s: %s\n",
					pattern, err.Error())
				continue
			}

			if matches {
				return f.dirPath + "/" + entry.Name(), true
			}
		}
	}
	return "", false
}

// IsPathDirectory checks if a given relative path is a directory.
//
// Returns true if it's a directory, or an error if the path can't be
// accessed.
func (f *Finder) IsPathDirectory(path string) (bool, error) {
	path = filepath.Join(f.dirPath, path)

	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to get stat %s: %w", path, err)
	}

	return info.IsDir(), nil
}