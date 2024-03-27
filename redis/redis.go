// A simple Dagger module to spawn and manage a Redis server.
//
// This module is designed for development and CI purposes only, do not use it to host a production server.
// The module implements a server and a client that can work together.
//
// You can use it to run a local Redis server but also in your CI pipeline to test your application
// with integrations tests.

package main

type Redis struct {
	Port     int
	Version  string
	Password *Secret
	Cache    bool
}

func New(
	// The port to use for the Redis server.
	//
	//+optional
	//+default=6379
	port int,

	// The version of the Redis server to use.
	//
	//+optional
	//+default="7.2.4"
	version string,

	// The password to use for the Redis server.
	//
	//+optional
	password *Secret,

	// Enable data persistency by mounting a cache volume.
	//
	//+optional
	//+default=false
	cache bool,	
) *Redis {
	return &Redis{
		Port:     port,
		Version:  version,
		Password: password,
		Cache:    cache,
	}
}
