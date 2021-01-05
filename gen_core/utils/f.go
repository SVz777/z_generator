/**
 * @file    f.go
 * @author
 *  ___  _  _  ____
 * / __)( \/ )(_   )
 * \__ \ \  /  / /_
 * (___/  \/  (____)
 * (903943711@qq.com)
 * @date    2019-12-09
 * @desc
 */
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
