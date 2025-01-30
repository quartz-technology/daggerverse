package dockercompose

// Cache represents a cache volume inside a container.
type Cache struct {
	// name is the cache volume's name.
	name string
	// path is the location where the cache is mounted inside the container.
	path string
}

// Name returns cache volume's name.
func (c *Cache) Name() string {
	return c.name
}

// Path returns the path to mount inside the container.
func (c *Cache) Path() string {
	return c.path
}
