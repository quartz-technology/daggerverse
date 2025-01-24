package dockercompose

type Cache struct {
	name string
	path string
}

func (c *Cache) Name() string {
	return c.name
}

func (c *Cache) Path() string {
	return c.path
}