package utils

import "dagger.io/dockersdk/module/object"

// MergeObjectsMap merges multiple maps of objects into a single map.
func MergeObjectsMap(objects ...map[string]object.Object) map[string]object.Object {
	merged := make(map[string]object.Object)

	for _, objects := range objects {
		for name, object := range objects {
			merged[name] = object
		}
	}

	return merged
}