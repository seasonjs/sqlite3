//go:build windows

package deps

import _ "embed"

//go:embed windows/sqlite-abi.dll
var libSqlite []byte

var libName = "sqlite-*.dll"

func getDylib() []byte {
	return libSqlite
}
