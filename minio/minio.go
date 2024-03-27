// A simple Dagger module to spawn and manage a Minio server.
//
// This module is designed for development and CI purposes only, do not use it to host a production server.
// The module implements a server and a client that can work together.
//
// It's designed to be used in other dagger modules but also as standalone.

package main

type Minio struct {
	Version     string
	ServerPort  int
	ConsolePort int
	Username    *Secret
	Password    *Secret
	Cache       bool
}

func New(
	// version of Minio to use.
	version string,

	// port to listen to the server.
	serverPort int,

	// port to listen to the console.
	consolePort int,

	// +optional
	// username to use.
	username *Secret,

	// +optional
	// username to use.
	password *Secret,

	// +optional
	// cache enables long living storage on the server.
	cache bool,
) *Minio {
	return &Minio{
		Version:     version,
		ServerPort:  serverPort,
		ConsolePort: consolePort,
		Username:    username,
		Password:    password,
		Cache:       cache,
	}
}
