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

type PublishOpts struct {
	Token   *Secret `doc:"Secret token to register to npm registry"`
	Version string  `doc:"Version of the package"`
	Access  string  `doc:"Access permission of the package"`
	DryRun  bool
}

// Publish publishes the source code to npm registry.
func (n *Node) Publish(ctx context.Context, opts PublishOpts) (string, error) {
	token, err := opts.Token.Plaintext(ctx)
	if err != nil {
		return "", err
	}

	// Configure .npmrc
	npmrc := fmt.Sprintf(`//registry.npmjs.org/:_authToken=%s
registry=https://registry.npmjs.org/
always-auth=true`, token)

	publishCmd := []string{"publish"}
	if opts.DryRun {
		publishCmd = append(publishCmd, "--dry-run")
	}

	if opts.Access != "" {
		publishCmd = append(publishCmd, "--access", opts.Access)
	}

	// Set version and publish
	return n.Ctr.
		WithNewFile(".npmrc", ContainerWithNewFileOpts{
			Contents:    npmrc,
			Permissions: 0o600,
		}).
		WithExec([]string{"version", "--new-version", opts.Version}).
		WithExec(publishCmd).
		Stdout(ctx)
}

type InstallOpts struct {
	Pkg []string `doc:"Package to additionally install"`
}

// Install adds given package.
func (n *Node) Install(opts InstallOpts) *Node {
	cmd := append([]string{"install"}, opts.Pkg...)

	return n.WithContainer(n.Run(cmd))
}
