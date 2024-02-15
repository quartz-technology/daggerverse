package main

import "strings"

func excludeModules(modules, exclude []string) []string {
	paths := []string{}

	// Exclude directories the user specified
	for _, path := range modules {
		excluded := false

		for _, e := range exclude {
			if strings.HasPrefix(path, e) {
				excluded = true
			}
		}

		if !excluded {
			paths = append(paths, path)
		}
	}

	return paths
}