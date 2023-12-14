package sqlite3

import (
	"errors"
	"github.com/seasonjs/sqlite3/deps"
	"os"
	"time"
)

// SQLiteTimestampFormats is timestamp formats understood by both this module
// and SQLite.  The first format in the slice will be used when saving time
// values into the database. When parsing a string from a timestamp or datetime
// column, the formats are tried in order.
var SQLiteTimestampFormats = []string{
	// By default, store timestamps with whatever timezone they come with.
	// When parsed, they will be returned with the same timezone.
	"2006-01-02 15:04:05.999999999-07:00",
	"2006-01-02T15:04:05.999999999-07:00",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02T15:04:05.999999999",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04",
	"2006-01-02T15:04",
	"2006-01-02",
}

type SQLite3 struct {
	cSQLite    CSQLite
	isAutoLoad bool
	dylibPath  string
}

func NewSQLite3Auto() (*SQLite3, error) {
	library, err := deps.DumpSQLiteLibrary()
	if err != nil {
		return nil, err
	}
	sqlite, err := NewSQLite3(library)
	if err != nil {
		return nil, err
	}
	sqlite.isAutoLoad = true
	sqlite.dylibPath = library
	return sqlite, nil
}

func NewSQLite3(dylibPath string) (*SQLite3, error) {
	cSQLite, err := NewCSQLite(dylibPath)
	if err != nil {
		return nil, err
	}
	return &SQLite3{
		cSQLite: cSQLite,
	}, nil
}

type Connect struct {
	ctx     *CSQLiteCtx
	sqlite3 *SQLite3
}

func (s *SQLite3) Open(path string) (*Connect, error) {
	ctx := s.cSQLite.SQLiteInit(path)
	if s.cSQLite.SQLiteErrCode(ctx) != 0 {
		return nil, errors.New(s.cSQLite.SQLiteErrMsg(ctx))
	}
	return &Connect{ctx: ctx, sqlite3: s}, nil
}

func (s *SQLite3) Close() error {
	err := s.cSQLite.Close()
	if err != nil {
		return err
	}
	if s.isAutoLoad {
		err = os.Remove(s.dylibPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Connect) Exec(sql string) (string, error) {
	if err := checkConnect(c.ctx); err != nil {
		return "", err
	}
	res := c.sqlite3.cSQLite.SQLiteExec(c.ctx, sql)
	if c.sqlite3.cSQLite.SQLiteErrCode(c.ctx) != 0 {
		return "", errors.New(c.sqlite3.cSQLite.SQLiteErrMsg(c.ctx))
	}
	return res, nil
}

func (c *Connect) Close() error {
	if err := checkConnect(c.ctx); err != nil {
		return err
	}
	res := c.sqlite3.cSQLite.SQLiteClose(c.ctx)
	if res != 0 {
		return errors.New(c.sqlite3.cSQLite.SQLiteErrMsg(c.ctx))
	}
	return nil
}

type Stmt struct {
	sql     string
	cStmt   *CSQLiteStmt
	sqlite3 *SQLite3
}

func (c *Connect) Prepare(sql string) (*Stmt, error) {
	if err := checkConnect(c.ctx); err != nil {
		return nil, err
	}
	res := c.sqlite3.cSQLite.SQLitePrepare(c.ctx, sql)
	if c.sqlite3.cSQLite.SQLiteErrCode(c.ctx) != 0 {
		return nil, errors.New(c.sqlite3.cSQLite.SQLiteErrMsg(c.ctx))
	}
	return &Stmt{cStmt: res, sql: sql, sqlite3: c.sqlite3}, nil
}

type Col struct {
	Name  string
	Value string
}

type Result struct {
	Rows  []Col
	total int
}

func (s *Stmt) Exec(args ...any) ([]Result, error) {
	for i, arg := range args {
		if arg == nil {
			s.sqlite3.cSQLite.SQLiteBindNull(s.cStmt, i+1)
			continue
		}
		switch arg.(type) {
		case int8, int16, int32, int:
			s.sqlite3.cSQLite.SQLiteBindInt(s.cStmt, i+1, arg.(int))
		case int64, uint64:
			s.sqlite3.cSQLite.SQLiteBindInt64(s.cStmt, i+1, arg.(int64))
		case string:
			s.sqlite3.cSQLite.SQLiteBindText(s.cStmt, i+1, arg.(string))
		case float32, float64:
			s.sqlite3.cSQLite.SQLiteBindDouble(s.cStmt, i+1, arg.(float64))
		case time.Time:
			s.sqlite3.cSQLite.SQLiteBindText(s.cStmt, i+1, arg.(time.Time).Format(SQLiteTimestampFormats[0]))
		default:
			return nil, errors.New("not support other type yet")
		}
	}

	res := make([]Result, 0)
	for {
		rc := s.sqlite3.cSQLite.SQLiteStep(s.cStmt)
		if rc == 101 {
			return res, nil
		}
		//res = append(res, Result{
		//	Name:  s.sqlite3.cSQLite.SQLiteColumnName(s.cStmt),
		//	Value: "",
		//})
	}

}

func checkConnect(ctx *CSQLiteCtx) error {
	if ctx.ctx == 0 {
		return errors.New("connect ctx is nil, you must call Open to start a connect")
	}
	return nil
}
