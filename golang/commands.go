package main

// Exec returns the container with the given command set.
func (g *Golang) Exec(args []string) *Container {
	return g.Ctr.WithExec(args)
}

// Run returns the container with the run command executed.
func (g *Golang) Run(file string) *Container {
	return g.Exec([]string{"run", file})
}

// Test returns the container with the test executed.
func (g *Golang) Test() *Container {
	return g.Exec([]string{"test"})
}

type BuildOpts struct {
	Output string `doc:"Path to write the built binary"`
}

// Build returns the container with the built artifact.
func (g *Golang) Build(
	// Path to write the built binary
	output Optional[string],
) *Container {
	cmd := []string{"build"}

	if output, ok := output.Get(); ok {
		cmd = append(cmd, "-o", output)
	}

	return g.Exec(cmd)
}

// Download installs dependencies written in the source directory.
func (g *Golang) Download() *Golang {
	return g.WithContainer(g.Exec([]string{"mod", "download"}))
}

// Get returns the container with given packages downloaded.
func (g *Golang) Get(packages []string) *Golang {
	cmd := append([]string{"get", "-u"}, packages...)

	return g.WithContainer(g.Exec(cmd))
}

// Install returns the container with given packages installed.
func (g *Golang) Install(packages ...string) *Golang {
	cmd := append([]string{"install"}, packages...)

	return g.WithContainer(g.Exec(cmd))
}
