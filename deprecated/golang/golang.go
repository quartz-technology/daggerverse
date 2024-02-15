package main

import "fmt"

type Golang struct {
	Ctr *Container
}

// WithVersion returns Golang container with given image version.
//
// The container is configured with cache for packages and build artifacts.
// The default entrypoint is set to `go`.
func (g *Golang) WithVersion(version string) *Golang {
	if g.Ctr == nil {
		g.Ctr = dag.Container()
	}

	g.Ctr = g.Ctr.
		From(fmt.Sprintf("golang:%s", version)).
		WithEntrypoint([]string{"go"}).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("gomod")).
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("gobuild"))

	return g
}

// WithContainer returns Golang container set with the given container.
func (g *Golang) WithContainer(ctr *Container) *Golang {
	g.Ctr = ctr

	return g
}

// Container returns Golang container.
func (g *Golang) Container() *Container {
	return g.Ctr
}

// WithSource returns the Golang container with given source mounted to `/src`.
func (g *Golang) WithSource(source *Directory) *Golang {
	workdir := "/src"

	g.Ctr = g.Ctr.
		WithWorkdir(workdir).
		WithMountedDirectory(workdir, source)

	return g
}
