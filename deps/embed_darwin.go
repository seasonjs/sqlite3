//go:build darwin

package deps

import _ "embed"

//go:embed darwin/libsqlite-abi.dylib
var libSqlite []byte

var libName = "libsqlite-*.dylib"

func getDylib() []byte {
	return libSqlite
}
