package main

type Redis struct {
	Port     int
	Version  string
	Password *Secret
}

func (r *Redis) WithPort(port int) *Redis {
	r.Port = port

	return r
}

func (r *Redis) WithVersion(version string) *Redis {
	r.Version = version

	return r
}

func (r *Redis) WithPassword(password *Secret) *Redis {
	r.Password = password

	return r
}
