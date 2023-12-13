package sqlite3

import "github.com/ebitengine/purego"

type CSQLiteCtx struct {
	ctx uintptr
}

type CSQLite interface {
	SQLiteInit(path string) *CSQLiteCtx
	SQLiteExec(ctx *CSQLiteCtx, sql string) string
	SQLiteClose(ctx *CSQLiteCtx) int
	SQLiteErrMsg(ctx *CSQLiteCtx) string
	SQLiteErrCode(ctx *CSQLiteCtx) int
	Close() error
}

type CSQLiteImpl struct {
	libSQLite     uintptr
	sqliteInit    func(path string) uintptr
	sqliteExec    func(db uintptr, sql string) string
	sqliteClose   func(db uintptr) int
	sqliteErrMsg  func(db uintptr) string
	sqliteErrCode func(db uintptr) int
}

func NewCSQLite(libraryPath string) (CSQLite, error) {
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

	return &CSQLiteImpl{
		libSQLite:     library,
		sqliteInit:    sqliteInit,
		sqliteExec:    sqliteExec,
		sqliteClose:   sqliteClose,
		sqliteErrMsg:  sqliteErrMsg,
		sqliteErrCode: sqliteErrCode,
	}, nil
}

func (c *CSQLiteImpl) SQLiteInit(path string) *CSQLiteCtx {
	ctx := c.sqliteInit(path)

	return &CSQLiteCtx{
		ctx,
	}
}

func (c *CSQLiteImpl) SQLiteExec(ctx *CSQLiteCtx, sql string) string {
	return c.sqliteExec(ctx.ctx, sql)
}

func (c *CSQLiteImpl) SQLiteClose(ctx *CSQLiteCtx) int {
	return c.sqliteClose(ctx.ctx)
}

func (c *CSQLiteImpl) SQLiteErrMsg(ctx *CSQLiteCtx) string {
	return c.sqliteErrMsg(ctx.ctx)
}

func (c *CSQLiteImpl) SQLiteErrCode(ctx *CSQLiteCtx) int {
	return c.sqliteErrCode(ctx.ctx)
}

func (c *CSQLiteImpl) Close() error {
	return closeLibrary(c.libSQLite)
}
