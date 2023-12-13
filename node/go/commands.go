package main

import (
	"context"
	"fmt"
)

// Run returns the container with the command set.
func (n *Node) Run(args []string) *Container {
	return n.Ctr.WithExec(args)
}

// Start returns the output of the code executed.
func (n *Node) Start(ctx context.Context) (string, error) {
	return n.Run([]string{"run", "start"}).Stdout(ctx)
}

// Lint returns the output of the lint command.
func (n *Node) Lint(ctx context.Context) (string, error) {
	return n.Run([]string{"run", "lint"}).Stdout(ctx)
}

// Test returns the result of the test executed.
func (n *Node) Test(ctx context.Context) (string, error) {
	return n.Run([]string{"run", "test"}).Stdout(ctx)
}

// Build returns the Container with the source built.
func (n *Node) Build() *Node {
	n.Ctr = n.Run([]string{"run", "build"})

	return n
}

// Publish publishes the source code to npm registry.
func (n *Node) Publish(
	ctx context.Context,
	// secret token to register to npm registry
	tokenSecret *Secret,
	// version of the package
	version string,
	// access permission of the package
	access Optional[string],
	// whether to do a dry run instead of an actual publish
	dryRun Optional[bool],
) (string, error) {
	token, err := tokenSecret.Plaintext(ctx)
	if err != nil {
		return "", err
	}

	// Configure .npmrc
	npmrc := fmt.Sprintf(`//registry.npmjs.org/:_authToken=%s
registry=https://registry.npmjs.org/
always-auth=true`, token)

	publishCmd := []string{"publish"}
	if dryRun.GetOr(false) {
		publishCmd = append(publishCmd, "--dry-run")
	}

	if access, ok := access.Get(); ok {
		publishCmd = append(publishCmd, "--access", access)
	}

	// Set version and publish
	return n.Ctr.
		WithNewFile(".npmrc", ContainerWithNewFileOpts{
			Contents:    npmrc,
			Permissions: 0o600,
		}).
		WithExec([]string{"version", "--new-version", version}).
		WithExec(publishCmd).
		Stdout(ctx)
}

// Install adds given package.
func (n *Node) Install(
	// packages to additionally install
	pkgs ...string,
) *Node {
	cmd := append([]string{"install"}, pkgs...)

	return n.WithContainer(n.Run(cmd))
}
