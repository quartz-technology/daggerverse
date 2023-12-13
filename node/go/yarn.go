package main

// WithYarn returns Node container with yarn configured as package manager.
func (n *Node) WithYarn() *Node {
	n.Ctr = n.Ctr.
		WithEntrypoint([]string{"yarn"}).
		WithMountedCache("/usr/local/share/.cache/yarn", dag.CacheVolume("yarn-cache"))

	return n
}
