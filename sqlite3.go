package sqlite3

import (
	"errors"
	"github.com/seasonjs/sqlite3/deps"
	"os"
)

type Sqlite3 struct {
	cSqlite    CSqlite
	isAutoLoad bool
	dylibPath  string
}

func NewSqlite3Auto() (*Sqlite3, error) {
	library, err := deps.DumpSqliteLibrary()
	if err != nil {
		return nil, err
	}
	sqlite, err := NewSqlite3(library)
	if err != nil {
		return nil, err
	}
	sqlite.isAutoLoad = true
	sqlite.dylibPath = library
	return sqlite, nil
}

func NewSqlite3(dylibPath string) (*Sqlite3, error) {
	cSqlite, err := NewCSqlite(dylibPath)
	if err != nil {
		return nil, err
	}
	return &Sqlite3{
		cSqlite: cSqlite,
	}, nil
}

type Connect struct {
	ctx     *CSqliteCtx
	sqlite3 *Sqlite3
}

func (s *Sqlite3) Open(path string) (*Connect, error) {
	ctx := s.cSqlite.SqliteInit(path)
	if s.cSqlite.SqliteErrCode(ctx) != 0 {
		return nil, errors.New(s.cSqlite.SqliteErrMsg(ctx))
	}
	return &Connect{ctx: ctx, sqlite3: s}, nil
}

func (s *Sqlite3) Close() error {
	err := s.cSqlite.Close()
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
	res := c.sqlite3.cSqlite.SqliteExec(c.ctx, sql)
	if c.sqlite3.cSqlite.SqliteErrCode(c.ctx) != 0 {
		return "", errors.New(c.sqlite3.cSqlite.SqliteErrMsg(c.ctx))
	}
	return res, nil
}

func (c *Connect) Close() error {
	if err := checkConnect(c.ctx); err != nil {
		return err
	}
	res := c.sqlite3.cSqlite.SqliteClose(c.ctx)
	if res != 0 {
		return errors.New(c.sqlite3.cSqlite.SqliteErrMsg(c.ctx))
	}
	return nil
}

func checkConnect(ctx *CSqliteCtx) error {
	if ctx.ctx == 0 {
		return errors.New("connect ctx is nil, you must call Open to start a connect")
	}
	return nil
}
