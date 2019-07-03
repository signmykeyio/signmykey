package common

import (
	"strings"
)

// TransformCase returns principals list in lower or upper case
func TransformCase(transform string, list []string) []string {
	principals := []string{}

	if transform == "lower" {
		for _, str := range list {
			principals = append(principals, strings.ToLower(str))
		}

		return principals
	}

	if transform == "upper" {
		for _, str := range list {
			principals = append(principals, strings.ToUpper(str))
		}

		return principals
	}

	return list
}
