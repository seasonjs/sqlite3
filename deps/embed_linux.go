//go:build linux

package deps

import _ "embed"

//go:embed linux/sqlite-abi.so
var libSqlite []byte

var libName = "sqlite-*.so"

func getDylib() []byte {
	return libSqlite
}
