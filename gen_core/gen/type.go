/**
 * @file    type.go
 * @author
 *  ___  _  _  ____
 * / __)( \/ )(_   )
 * \__ \ \  /  / /_
 * (___/  \/  (____)
 * (903943711@qq.com)
 * @date    2020-01-04
 * @desc
 */
package gen

import "bytes"

type IGen interface {
	Gen() *bytes.Buffer
}
