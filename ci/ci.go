package main

type Ci struct{}

func source() *Directory {
	return dag.
		Host().
		Directory(".", HostDirectoryOpts{
			Include: []string{"**/*.go", "**/go.mod", "**/go.sum", ".golangci.yml"},
		})
}

func daggerRepository() *Directory {
	return dag.
		Git("github.com/dagger/dagger").
		Branch("main").
		Tree()
}
