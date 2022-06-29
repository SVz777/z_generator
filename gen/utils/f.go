package utils

import (
	"strings"
)

func IsValueType(typeName string) bool {
	if strings.HasPrefix(typeName, "map") || strings.HasPrefix(typeName, "[]") {
		return false
	}
	return true
}
