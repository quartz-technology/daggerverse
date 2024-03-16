package main

import (
	"context"
	"fmt"
	"strings"
)

// CLI is an abstraction to the Dagger CLI so you can use it inside a container
// for high level management command.
type CLI struct {
	// The Dagger CLI container
	Container *Container
}

// CLI installs the Dagger CLI in a container and returns it.
//
// Make sure to set the `version` to the desired Dagger version.
// You may also need to enable `Experimental Privileged Nesting` in your container to make it works.
//
// Example usage:
//
//	dagger call cli --version=0.10.2 with-exec --args="version" stdout
func (d *DaggerverseCockpit) CLI(
	ctx context.Context,

	//+optional
	//+default="0.10.2"
	version string,
) *CLI {
	return &CLI{
		Container: dag.
			Container().
			From("alpine:3.19.1").
			WithExec([]string{"apk", "add", "curl"}).
			WithExec([]string{
				"sh", "-c",
				fmt.Sprintf("curl -L https://dl.dagger.io/dagger/install.sh | %s sh", fmt.Sprintf("DAGGER_VERSION=%s", version)),
			}).
			WithDefaultTerminalCmd([]string{"sh"}, ContainerWithDefaultTerminalCmdOpts{
				ExperimentalPrivilegedNesting: true,
				InsecureRootCapabilities:      true,
			}).
			WithEntrypoint([]string{"/bin/dagger"}),
	}
}

// Publish executes the publish command to upload the module to the Daggerverse.
// This function returns the URL of the published module.
// Example usage:
//  dagger call cli publish --repository=. 
//
// You could also use it directly in your code:
//  repository := // ... your module repository fetched from arguments or git
//  url, err := dag.DaggerCockpit().CLI().Publish(ctx, repository, ".")
//  if err !=  nil {...} // Handle error
//  fmt.Println(url)
func (c *CLI) Publish(
	ctx context.Context,

	// The repository to use the Dagger CLI on.
	// Dagger expect the `.git` directories to be inside this directory.
	// Specify a subpath if your module is located in a child directory.
	repository *Directory,

	// The path to the module to publish
	//+optional
	//+default="."
	path string,
) (string, error) {
	workdir := "/module"
	
	out, err := c.
	  Container.
	  WithWorkdir(workdir).
	  WithDirectory(workdir, repository).
	  WithExec([]string{"publish", "-m", path}, ContainerWithExecOpts{ExperimentalPrivilegedNesting: true}).
	  Stdout(ctx)

	if err != nil {
		return "", fmt.Errorf("could not publish the module: %w", err)
	}

	// Check the logs to get the URL of the module
	logs := strings.Split(out, "\n")
	return strings.TrimSpace(logs[len(logs)-2]), nil
}

// Develop executes the develop command to start a development environment for the module
// and returns its content.
//
// Example usage:
//  dagger call cli develop --module=. -o .
func (c *CLI) Develop(
	ctx context.Context,

	// The module to use the Dagger CLI on.
	module *Directory,
) *Directory {
	workdir := "/module"
	
	generatedDir := c.
	  Container.
	  WithWorkdir(workdir).
	  WithDirectory(workdir, module).
	  WithExec([]string{"develop"}, ContainerWithExecOpts{ExperimentalPrivilegedNesting: true}).
	  Directory(workdir)

	return module.Diff(generatedDir)
}