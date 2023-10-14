package main

import (
	"context"
	"fmt"
	"os"
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
	Version string  `doc:"version of the package"`
	Access  string  `doc:"access permission of the package"`
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
	if err := os.WriteFile(".npmrc", []byte(npmrc), 0o600); err != nil {
		return "", err
	}

	// Set version and publish
	return n.Ctr.
		WithExec([]string{"version", opts.Version}).
		WithExec([]string{"publish", "--access", opts.Access}).
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
