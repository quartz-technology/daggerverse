package main

import (
	"fmt"
)

type Node struct {
	Ctr *Container
}

// WithVersion returns Node container with given image version.
func (n *Node) WithVersion(version string) *Node {
	if n.Ctr == nil {
		n.Ctr = dag.Container()
	}

	n.Ctr = n.Ctr.
		From(fmt.Sprintf("node:%s", version)).
		WithEntrypoint([]string{"node"})

	return n
}

// WithContainer returns Node container with the given container.
func (n *Node) WithContainer(ctr *Container) *Node {
	n.Ctr = ctr

	return n
}

// Container returns Node container.
func (n *Node) Container() *Container {
	return n.Ctr
}

// WithSource returns the Node container with source and cache set in it.
func (n *Node) WithSource(source *Directory) *Node {
	workdir := "/src"

	n.Ctr = n.Ctr.
		WithWorkdir("/src").
		WithMountedDirectory("/src", source).
		WithMountedCache(
			fmt.Sprintf("%s/node_modules", workdir),
			dag.CacheVolume("node-modules"))

	return n
}
