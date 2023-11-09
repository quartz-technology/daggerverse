package main

import "context"

// Publish execute the publish command in the current workdir.
//
// Use this command to publish a module to dagger-verse
func (c *CLI) Publish(
	ctx context.Context,

	// Specify a path to the module located in a subdirectory of the workdir
	// (e.g., "./my-cool-module")
	path Optional[string],
) (string, error) {
	module := path.GetOr(".")

	return c.
		Container().
		WithExec(
			[]string{"mod", "publish", "-m", module, "-f"},
			ContainerWithExecOpts{ExperimentalPrivilegedNesting: true},
		).
		Stdout(ctx)
}

// Call returns the result of the call command.
func (c *CLI) Call(
	ctx context.Context,

	// Command to call (e.g., ["integration-test", "run"])
	command []string,

	// Arguments to add to the command (e.g., ["--target", "foo"])
	args ...string,
) (string, error) {
	args = append([]string{"call"}, append(command, args...)...)

	return c.
		Container().
		WithExec(
			args,
			ContainerWithExecOpts{ExperimentalPrivilegedNesting: true}).
		Stdout(ctx)
}

// Run returns the result of the run command.
func (c *CLI) Run(
	ctx context.Context,

	// Arguments to pass to run (e.g., ["echo", "hello world"])
	args ...string,
) (string, error) {
	args = append([]string{"run"}, args...)

	return c.
		Container().
		WithExec(
			args,
			ContainerWithExecOpts{ExperimentalPrivilegedNesting: true}).
		Stdout(ctx)
}

// Query returns the result of a graphQL query processed by dagger.
func (c *CLI) Query(ctx context.Context, query string) (string, error) {
	args := append([]string{"query"}, query)

	return c.
		Container().
		WithExec(
			args,
			ContainerWithExecOpts{ExperimentalPrivilegedNesting: true}).
		Stdout(ctx)
}
