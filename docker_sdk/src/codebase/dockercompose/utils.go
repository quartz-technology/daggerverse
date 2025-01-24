package dockercompose

import (
	"fmt"
	"strings"
)

// For some reason, docker-compose prefix all relative paths with /scratch.
// Because we work with the directory itself, we need to remove it.
func trimHostPath(hostPath string) string {
	if hostPath == "/scratch" {
		return "."
	}

	return fmt.Sprintf("./%s", strings.TrimPrefix(hostPath, "/scratch/"))
}