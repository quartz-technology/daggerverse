package dockercompose

// SourceImage represents a source based on an image reference.
type SourceImage struct {
	Ref string // Image reference.
}

// SourceDockerfile represents a source configured using a Dockerfile.
type SourceDockerfile struct {
	// Build context path.
	Context string
	// Dockerfile path.
	Dockerfile string

	// Arguments to pass during the build process.
	BuildArgs map[string]*string

	// Specific target in a multi-stage Dockerfile.
	Target *string
}

// SourceType defines the type of source used, either an image or a Dockerfile.
type SourceType string

const (
	// SourceTypeImage indicates that the source is an image reference.
	SourceTypeImage SourceType = "image"
	// SourceTypeDockerfile indicates that the source is described by a Dockerfile.
	SourceTypeDockerfile SourceType = "dockerfile"
)

// Source encapsulates a source definition, identifying whether it uses an
// image or a Dockerfile, alongside the necessary configuration details.
type Source struct {
	Type       SourceType
	Image      *SourceImage
	Dockerfile *SourceDockerfile
}
