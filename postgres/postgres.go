package main

type Postgres struct {
	User        *Secret
	Password    *Secret
	Name        string
	Port        int
	Version     string
	ConfigFile  *File
	InitScripts *Directory
	Cache       bool
}

func New(
	// User to access the database
	// This value will be set as a secret variable in the container: `POSTGRES_USER`
	user *Secret,

	// Password to access the database.
	// This value will be set as a secret variable in the container: `POSTGRES_PASSWORD`
	password *Secret,

	// Port to expose the database on.
	//+default=5432
	dbPort int,

	// Name of the database to create on start.
	// If it's not set, the database's name will be the user's one.
	//
	//+optional
	dbName string,

	// Enable data persistency by using a volume.
	//
	//+optional
	//+default=false
	cache bool,

	// PostgreSQL version to pull from the registry.
	//
	//+optional
	//+default="16.2"
	version string,

	// Extra configuration file to add to the postgres database.
	// This file will be copied in the container to `/usr/share/postgresql/postgresql.conf`
	//
	//+optional
	configFile *File,

	// Scripts to execute when the database is started.
	// The files will be copied to `/docker-entrypoint-initdb.d`
	//
	//+optional
	initScript *Directory,
) *Postgres {
	return &Postgres{
		User:        user,
		Password:    password,
		Name:        dbName,
		Port:        dbPort,
		Version:     version,
		ConfigFile:  configFile,
		InitScripts: initScript,
		Cache:       cache,
	}
}
