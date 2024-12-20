package utils

import (
	"dagger.io/dagger"
	"dagger.io/dagger/dag"

	"encoding/json"
	"fmt"
)

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


