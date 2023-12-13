package sqlite3

import (
	"errors"
	"github.com/seasonjs/sqlite3/deps"
	"os"
)

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

func checkConnect(ctx *CSQLiteCtx) error {
	if ctx.ctx == 0 {
		return errors.New("connect ctx is nil, you must call Open to start a connect")
	}
	return nil
}
