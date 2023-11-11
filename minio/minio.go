package main

type Minio struct {
	Version     string
	ServerPort  int
	ConsolePort int
	Username    *Secret
	Password    *Secret
	Cache       bool
}

// WithVersion sets a version for Minio.
func (m *Minio) WithVersion(version string) *Minio {
	m.Version = version

	return m
}

// WithConsolePort sets a port to listen to the console.
func (m *Minio) WithConsolePort(port int) *Minio {
	m.ConsolePort = port

	return m
}

// WithServerPort sets a port to listen to the server.
func (m *Minio) WithServerPort(port int) *Minio {
	m.ServerPort = port

	return m
}

// WithCredential sets access and secret key in the CLI.
func (m *Minio) WithCredential(username *Secret, password *Secret) *Minio {
	m.Username = username
	m.Password = password

	return m
}

// WithCache enables long living storage on the server.
func (m *Minio) WithCache(cache bool) *Minio {
	m.Cache = cache

	return m
}
