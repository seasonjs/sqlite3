//go:build linux

package deps

import _ "embed"

//go:embed linux/libsqlite-abi.so
var libSqlite []byte

var libName = "libsqlite-*.so"

func getDylib() []byte {
	return libSqlite
}
