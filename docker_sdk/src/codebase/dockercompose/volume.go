package dockercompose

import "path"

// Volume represents an abstraction of a Docker Compose volume with its origin
// and target paths.
type Volume struct {
	// origin specifies the path of the volume on the host system.
	origin string
	// target specifies the desired mount path inside the container.
	target string
	// isDir indicates whether the volume is a directory.
	isDir bool
}

// Name returns the base name of the origin path.
func (v *Volume) Name() string {
	base := path.Base(v.origin)
	if base == "." {
		// It's a special case where path is '.', so we set name it to current-directory
		return "current-directory"
	}

	return base
}

// Origin returns the origin path of the volume.
func (v *Volume) Origin() string {
	return v.origin
}

// Target returns the target path of the volume to mount inside the container.
func (v *Volume) Target() string {
	return v.target
}

// IsDir returns true if the volume is a directory.
func (v *Volume) IsDir() bool {
	return v.isDir
}