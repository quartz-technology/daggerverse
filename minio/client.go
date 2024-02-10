package main

import (
	"context"
	"fmt"
)

type MC struct {
	Ctr *Container
}

// MC returns a Minio Client.
func (m *Minio) MC() *MC {
	ctr := dag.
		Container().
		From(fmt.Sprintf("minio/mc:%s", m.Version))

	return &MC{
		Ctr: ctr,
	}
}

// MCFromContainer use the given container as Minio client.
//
// This is useful to configure your own minio client or execute special
// operation.
func (m *Minio) MCFromContainer(ctr *Container) *MC {
	return &MC{
		Ctr: ctr,
	}
}

// Container returns the Minio Container.
func (c *MC) Container() *Container {
	return c.Ctr
}

// AliasSet adds an alias to the minio client.
func (c *MC) AliasSet(alias string, host string, username *Secret, password *Secret) *MC {
	c.Ctr = c.Ctr.
		WithSecretVariable("MINIO_ACCESS_KEY", username).
		WithSecretVariable("MINIO_SECRET_KEY", password).
		WithExec(
			[]string{"/bin/sh", "-c", "mc", "alias", "set", alias, host, "$MINIO_ACCESS_KEY", "$MINIO_SECRET_KEY"},
			ContainerWithExecOpts{SkipEntrypoint: true},
		)

	return c
}

// AliasRemove deletes an alias from the minio client.
func (c *MC) AliasRemove(alias string) *MC {
	c.Ctr = c.Ctr.
		WithExec([]string{"remove", alias})

	return c
}

// List returns a list of all object gave at the target.
func (c *MC) List(ctx context.Context, target string) (string, error) {
	return c.Ctr.
		WithExec([]string{"ls", target}).
		Stdout(ctx)
}

// MakeBucket creates a bucket in the given target at the given path.
func (c *MC) MakeBucket(target string, path string) *MC {
	c.Ctr = c.Ctr.
		WithExec([]string{"mb", fmt.Sprintf("%s/%s", target, path)})

	return c
}

// CopyFile adds the given file to the path given in the target.
func (c *MC) CopyFile(file *File, target string, path string) *MC {
	c.Ctr = c.Ctr.
		WithMountedFile("/file.txt", file).
		WithExec([]string{"cp", "/file.txt", fmt.Sprintf("%s/%s", target, path)})

	return c
}

// CopyDir adds the given directory to the path given in the target.
func (c *MC) CopyDir(dir *Directory, target string, path string) *MC {
	c.Ctr = c.Ctr.
		WithMountedDirectory("/dir", dir).
		WithExec([]string{"cp", "/dir", fmt.Sprintf("%s/%s", target, path)})

	return c
}

// Exec returns the result of the given command.
func (c *MC) Exec(ctx context.Context, command ...string) (string, error) {
	return c.Ctr.
		WithExec(command).
		Stdout(ctx)
}
