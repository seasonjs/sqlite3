package sqlite3

import (
	"fmt"
	"os"
	"runtime"
	"testing"
)

func getLibrary() string {
	switch runtime.GOOS {
	case "darwin":
		return "./deps/darwin/libsqlite-abi.dylib"
	case "linux":
		return "./deps/linux/libsqlite-abi.so"
	case "windows":
		return "./deps/windows/sqlite-abi.dll"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}
}

func TestNewCSQLite(t *testing.T) {
	sqlite, err := NewCSQLite(getLibrary())
	if err != nil {
		t.Error(err)
		return
	}
	defer sqlite.Close()
	ctx := sqlite.SQLiteInit("./tmp.db")

	errCode := sqlite.SQLiteErrCode(ctx)
	t.Log(errCode)
	errMsg := sqlite.SQLiteErrMsg(ctx)
	t.Log(errMsg)

	res := sqlite.SQLiteExec(ctx, `
				CREATE TABLE Users (    
					id INT PRIMARY KEY,    
					name VARCHAR(100),    
					age INT,    
					email VARCHAR(100),   
					created_at DATETIME
                );
	`)
	errCode = sqlite.SQLiteErrCode(ctx)
	t.Log(errCode)
	errMsg = sqlite.SQLiteErrMsg(ctx)
	t.Log(errMsg)
	t.Log(res)

	sqlite.SQLiteClose(ctx)

	err = os.Remove("./tmp.db")
	if err != nil {
		t.Error(err)
		return
	}

}
