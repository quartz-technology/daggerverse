package main

// WithNpm returns Node container with npm configured as package manager.
func (n *Node) WithNpm() *Node {
	n.Ctr = n.Ctr.
		WithEntrypoint([]string{"npm"}).
		WithMountedCache("/root/.npm", dag.CacheVolume("npm-cache"))

	return n
}
