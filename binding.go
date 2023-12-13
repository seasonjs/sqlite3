package sqlite3

import "github.com/ebitengine/purego"

type CSqliteCtx struct {
	ctx uintptr
}

type CSqlite interface {
	SqliteInit(path string) *CSqliteCtx
	SqliteExec(ctx *CSqliteCtx, sql string) string
	SqliteClose(ctx *CSqliteCtx) int
	SqliteErrMsg(ctx *CSqliteCtx) string
	SqliteErrCode(ctx *CSqliteCtx) int
	Close() error
}

type CSqliteImpl struct {
	libSqlite     uintptr
	sqliteInit    func(path string) uintptr
	sqliteExec    func(db uintptr, sql string) string
	sqliteClose   func(db uintptr) int
	sqliteErrMsg  func(db uintptr) string
	sqliteErrCode func(db uintptr) int
}

func NewCSqlite(libraryPath string) (CSqlite, error) {
	library, err := openLibrary(libraryPath)
	if err != nil {
		return nil, err
	}
	var (
		sqliteInit    func(path string) uintptr
		sqliteExec    func(db uintptr, sql string) string
		sqliteClose   func(db uintptr) int
		sqliteErrMsg  func(db uintptr) string
		sqliteErrCode func(db uintptr) int
	)

	purego.RegisterLibFunc(&sqliteInit, library, "sqlite3_abi_init")
	purego.RegisterLibFunc(&sqliteExec, library, "sqlite3_abi_exec")
	purego.RegisterLibFunc(&sqliteClose, library, "sqlite3_abi_close")
	purego.RegisterLibFunc(&sqliteErrMsg, library, "sqlite3_abi_errmsg")
	purego.RegisterLibFunc(&sqliteErrCode, library, "sqlite3_abi_errcode")

	return &CSqliteImpl{
		libSqlite:     library,
		sqliteInit:    sqliteInit,
		sqliteExec:    sqliteExec,
		sqliteClose:   sqliteClose,
		sqliteErrMsg:  sqliteErrMsg,
		sqliteErrCode: sqliteErrCode,
	}, nil
}

func (c *CSqliteImpl) SqliteInit(path string) *CSqliteCtx {
	ctx := c.sqliteInit(path)

	return &CSqliteCtx{
		ctx,
	}
}

func (c *CSqliteImpl) SqliteExec(ctx *CSqliteCtx, sql string) string {
	return c.sqliteExec(ctx.ctx, sql)
}

func (c *CSqliteImpl) SqliteClose(ctx *CSqliteCtx) int {
	return c.sqliteClose(ctx.ctx)
}

func (c *CSqliteImpl) SqliteErrMsg(ctx *CSqliteCtx) string {
	return c.sqliteErrMsg(ctx.ctx)
}

func (c *CSqliteImpl) SqliteErrCode(ctx *CSqliteCtx) int {
	return c.sqliteErrCode(ctx.ctx)
}

func (c *CSqliteImpl) Close() error {
	return closeLibrary(c.libSqlite)
}
