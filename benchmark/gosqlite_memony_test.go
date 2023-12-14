package benchmark

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/seasonjs/sqlite3"
	"gorm.io/gorm"
	"testing"
)

//goos: windows
//goarch: amd64
//pkg: github.com/seasonjs/sqlite3/benchmark
//cpu: Intel(R) Core(TM) i7-9700 CPU @ 3.00GHz
//BenchmarkGorm
//BenchmarkGorm-8           116439              9887 ns/op
//PASS

func BenchmarkGorm(b *testing.B) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		b.Error(err)
		return
	}
	db.Exec(`
		CREATE TABLE Users (    
			id INT PRIMARY KEY,    
			name VARCHAR(100),    
			age INT,    
			email VARCHAR(100),   
			created_at DATETIME
		);
	`)
	for i := 0; i < b.N; i++ {
		db.Exec(`insert into Users(id)  values (?)`, i)
	}
}

// goos: windows
// goarch: amd64
// pkg: github.com/seasonjs/sqlite3/benchmark
// cpu: Intel(R) Core(TM) i7-9700 CPU @ 3.00GHz
// BenchmarkDLL
// BenchmarkDLL-8            182738              6608 ns/op
// PASS
func BenchmarkDLL(b *testing.B) {
	instance, err := sqlite3.NewSQLite3Auto()
	db, err := instance.Open(":memory:")
	if err != nil {
		b.Error(err)
		return
	}
	_, err = db.Exec(`
		CREATE TABLE Users (    
			id INT PRIMARY KEY,    
			name VARCHAR(100),    
			age INT,    
			email VARCHAR(100),   
			created_at DATETIME
		);
	`)
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		_, err = db.Exec(fmt.Sprintf("insert into Users(id)  values ( %d )", i))
		if err != nil {
			b.Error(err)
			return
		}
	}
}
