package benchmark

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/seasonjs/sqlite3"
	"gorm.io/gorm"
	"testing"
)

// goos: windows
// goarch: amd64
// pkg: github.com/seasonjs/sqlite3/benchmark
// cpu: Intel(R) Core(TM) i7-9700 CPU @ 3.00GHz
// BenchmarkGormFile
// BenchmarkGormFile-8         249           5024135 ns/op
// BenchmarkGormFile-8         220           5071035 ns/op
// PASS
func BenchmarkGormFile(b *testing.B) {
	db, err := gorm.Open(sqlite.Open("./a.db"), &gorm.Config{})
	if err != nil {
		b.Error(err)
		return
	}
	db.Exec(`
		 CREATE TABLE IF NOT EXISTS Users(    
			id INT PRIMARY KEY,    
			name VARCHAR(100),    
			age INT,    
			email VARCHAR(100),   
			created_at DATETIME
		);
	`)
	for i := 0; i < b.N; i++ {
		db.Exec(`insert into Users(age)  values (?)`, i)
	}
}

// goos: windows
// goarch: amd64
// pkg: github.com/seasonjs/sqlite3/benchmark
// cpu: Intel(R) Core(TM) i7-9700 CPU @ 3.00GHz
// BenchmarkDLLFile
// BenchmarkDLLFile-8            207              5174840 ns/op
// PASS
func BenchmarkDLLFile(b *testing.B) {
	//var lock sync.Mutex
	instance, err := sqlite3.NewSQLite3Auto()
	db, err := instance.Open("./b.db")

	if err != nil {
		b.Error(err)
		return
	}

	_, err = db.Exec(`
		 CREATE TABLE IF NOT EXISTS Users(    
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
		//not safe and not performance best,but easy to use
		_, err = db.Exec(fmt.Sprintf("insert into Users(age)  values ( %d )", i))
		if err != nil {
			b.Error(err)
			return
		}
	}
}
