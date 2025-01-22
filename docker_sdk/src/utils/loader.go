package utils

import (
	"dagger.io/dagger"
	"dagger.io/dagger/dag"
	"dagger.io/dockersdk/module/object"

	"encoding/json"
	"fmt"
)

// LoadArgument is a generic function to load an argument from 
// the input map.
// If the argument is not set, it returns the zero value of the type.
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

func LoadDefaultValue(value interface{}) dagger.JSON {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return ""
	}

	return dagger.JSON(jsonValue)
}

func LoadDirectoryFromID(idPayload []byte) *dagger.Directory {
	var id dagger.DirectoryID

	err := json.Unmarshal(idPayload, &id)
	if err != nil {
		panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg", err))
	}

	return dag.LoadDirectoryFromID(id)
}

func LoadContainerFromID(idPayload []byte) *dagger.Container {
	var id dagger.ContainerID

	err := json.Unmarshal(idPayload, &id)
	if err != nil {
		panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg", err))
	}

	return dag.LoadContainerFromID(id)
}

func LoadFileFromID(idPayload []byte) *dagger.File {
	var id dagger.FileID

	err := json.Unmarshal(idPayload, &id)
	if err != nil {
		panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg", err))
	}

	return dag.LoadFileFromID(id)
}

func LoadSecretFromID(idPayload []byte) *dagger.Secret {
	var id dagger.SecretID

	err := json.Unmarshal(idPayload, &id)
	if err != nil {
		panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg", err))
	}

	return dag.LoadSecretFromID(id)
}


