package utils

import "strings"

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