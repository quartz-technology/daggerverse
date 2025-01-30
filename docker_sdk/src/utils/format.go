package utils

import "strings"

// FormatEnvVariableName formats the given environment variable name
//
// It converts the name to uppercase and replaces underscores with
// uppercase letters.
//
// For example, the input "FOO_BAR" will be converted to "FooBar".
//
// This exists because some environment variables names triggers
// errors when used as function arguments.
// For example: `FOO_X_BAR` will fails with an unset argument error
// but if we convert it to `FooXBar`, it will work.
func FormatEnvVariableName(name string) string {
	parts := strings.Split(name, "_")
	
	res := []string{}
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		res = append(res, strings.ToUpper(part[0:1]) + strings.ToLower(part[1:]))
	}

	return strings.Join(res, "")
}

func RemoveListDuplicates[T comparable](list []T) []T {
	keys := make(map[T]bool)
	list2 := []T{}

	for _, entry := range list {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list2 = append(list2, entry)
		}
	}
	return list2
}