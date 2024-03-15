package main

type Redis struct {
	Port     int
	Version  string
	Password *Secret
	Cache    bool
}

func NewRedis(
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
