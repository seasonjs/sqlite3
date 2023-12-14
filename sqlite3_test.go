package sqlite3_test

import (
	"github.com/seasonjs/sqlite3"
	"os"
	"testing"
)

func TestNewSQLite3Auto(t *testing.T) {
	sqlite, err := sqlite3.NewSQLite3Auto()
	if err != nil {
		t.Error(err)
		return
	}
	defer func(sqlite *sqlite3.SQLite3) {
		err := sqlite.Close()
		if err != nil {
			t.Error(err)
		}
	}(sqlite)

	conn, err := sqlite.Open("./tmp.db")
	if err != nil {
		t.Error(err)
		return
	}

	exec, err := conn.Exec(`
				CREATE TABLE Users (    
					id INT PRIMARY KEY,    
					name VARCHAR(100),    
					age INT,    
					email VARCHAR(100),   
					created_at DATETIME
                );
	`)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(exec)

	err = conn.Close()
	if err != nil {
		t.Error(err)
		return
	}

	err = os.Remove("./tmp.db")
	if err != nil {
		t.Error(err)
		return
	}

}
