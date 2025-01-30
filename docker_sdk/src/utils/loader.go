// This file contains utility functions for handling JSON with dagger types.
//
// These functions provide conversion tools that facilitate mapping JSON 
// payloads to specific dagger entity types and handling default values.

package utils

import (
	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/module/object"

	"encoding/json"
	"fmt"
)

// LoadArgument loads a typed argument from an input map.
// Returns the zero value of the type if the argument is not set or unmarshalling fails.
func LoadArgument[K any](name string, args object.InputArgs) K {
	var res K

	if args[name] == nil {
		return res
	}

	err := json.Unmarshal(args[name], &res)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal input arg %s: %w", name, err))
	}

	return res
}

// LoadDefaultValue marshals a value into dagger.JSON.
// Returns an empty string if marshalling fails.
func LoadDefaultValue(value interface{}) dagger.JSON {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return ""
	}

	return dagger.JSON(jsonValue)
}

// LoadDirectoryFromID converts a byte payload to a dagger.Directory.
// It unmarshals the payload to obtain the DirectoryID and uses it to load the directory.
func LoadDirectoryFromID(idPayload []byte) *dagger.Directory {
	var id dagger.DirectoryID

	err := json.Unmarshal(idPayload, &id)
	if err != nil {
		panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg", err))
	}

	return dag.LoadDirectoryFromID(id)
}

// LoadContainerFromID converts a byte payload to a dagger.Container.
// It unmarshals the payload to obtain the ContainerID and uses it to load the container.
func LoadContainerFromID(idPayload []byte) *dagger.Container {
	var id dagger.ContainerID

	err := json.Unmarshal(idPayload, &id)
	if err != nil {
		panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg", err))
	}

	return dag.LoadContainerFromID(id)
}

// LoadFileFromID converts a byte payload to a dagger.File.
// It unmarshals the payload to obtain the FileID and uses it to load the file.
func LoadFileFromID(idPayload []byte) *dagger.File {
	var id dagger.FileID

	err := json.Unmarshal(idPayload, &id)
	if err != nil {
		panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg", err))
	}

	return dag.LoadFileFromID(id)
}

// LoadSecretFromID converts a byte payload to a dagger.Secret.
// It unmarshals the payload to obtain the SecretID and uses it to load the secret.
func LoadSecretFromID(idPayload []byte) *dagger.Secret {
	var id dagger.SecretID

	err := json.Unmarshal(idPayload, &id)
	if err != nil {
		panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg", err))
	}

	return dag.LoadSecretFromID(id)
}


