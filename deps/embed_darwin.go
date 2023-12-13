//go:build darwin

package deps

import _ "embed"

//go:embed darwin/sqlite-abi.dylib
var libSqlite []byte

var libName = "sqlite-*.dylib"

func getDylib() []byte {
	return libSqlite
}
