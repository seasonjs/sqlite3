package sqlite3

import "github.com/ebitengine/purego"

type CSQLiteCtx struct {
	ctx uintptr
}

type CSQLiteStmt struct {
	stmt uintptr
}

type CSQLite interface {
	SQLiteInit(path string) *CSQLiteCtx
	SQLiteExec(ctx *CSQLiteCtx, sql string) string
	SQLiteClose(ctx *CSQLiteCtx) int
	SQLiteErrMsg(ctx *CSQLiteCtx) string
	SQLiteErrCode(ctx *CSQLiteCtx) int
	Close() error

	SQLitePrepare(ctx *CSQLiteCtx, zSql string) *CSQLiteStmt
	SQLiteFinalize(stmt *CSQLiteStmt) int
	SQLiteReset(stmt *CSQLiteStmt) int
	SQLiteStep(stmt *CSQLiteStmt) int

	SQLiteBindBlob(stmt *CSQLiteStmt, index int, zData []byte) int
	SQLiteBindBlob64(stmt *CSQLiteStmt, index int, zData []byte) int
	SQLiteBindDouble(stmt *CSQLiteStmt, index int, zData float64) int
	SQLiteBindInt(stmt *CSQLiteStmt, index int, zData int) int
	SQLiteBindInt64(stmt *CSQLiteStmt, index int, zData int64) int
	SQLiteBindNull(stmt *CSQLiteStmt, index int) int
	SQLiteBindText(stmt *CSQLiteStmt, index int, zData string) int

	SQLiteColumnName(stmt *CSQLiteStmt, index int) string
	SQLiteColumnText(stmt *CSQLiteStmt, index int) string
}

type CSQLiteImpl struct {
	libSQLite     uintptr
	sqliteInit    func(path string) uintptr
	sqliteExec    func(db uintptr, sql string) string
	sqliteClose   func(db uintptr) int
	sqliteErrMsg  func(db uintptr) string
	sqliteErrCode func(db uintptr) int

	sqlitePrepareV2 func(db uintptr, zSql string, nByte int) uintptr
	sqliteFinalize  func(stmt uintptr) int
	sqliteReset     func(stmt uintptr) int
	sqliteStep      func(stmt uintptr) int

	sqliteBindBlob   func(stmt uintptr, index int, zData *byte, nData int, callback uintptr) int
	sqliteBindBlob64 func(stmt uintptr, index int, zData *byte, nData uint64, callback uintptr) int
	sqliteBindDouble func(stmt uintptr, index int, zData float64) int
	sqliteBindInt    func(stmt uintptr, index int, zData int) int
	sqliteBindInt64  func(stmt uintptr, index int, zData int64) int
	sqliteBindNull   func(stmt uintptr, index int) int
	sqliteBindText   func(stmt uintptr, index int, zData string, nData int) int

	sqlite3ColumnName func(stmt uintptr, index int) string
	sqlite3ColumnText func(stmt uintptr, index int) string
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

		sqlitePrepareV2 func(db uintptr, zSql string, nByte int) uintptr
		sqliteFinalize  func(stmt uintptr) int
		sqliteReset     func(stmt uintptr) int
		sqliteStep      func(stmt uintptr) int

		sqliteBindBlob   func(stmt uintptr, index int, zData *byte, nData int, callback uintptr) int
		sqliteBindBlob64 func(stmt uintptr, index int, zData *byte, nData uint64, callback uintptr) int
		sqliteBindDouble func(stmt uintptr, index int, zData float64) int
		sqliteBindInt    func(stmt uintptr, index int, zData int) int
		sqliteBindInt64  func(stmt uintptr, index int, zData int64) int
		sqliteBindNull   func(stmt uintptr, index int) int
		sqliteBindText   func(stmt uintptr, index int, zData string, nData int) int

		sqlite3ColumnName func(stmt uintptr, index int) string
		sqlite3ColumnText func(stmt uintptr, index int) string
	)

	purego.RegisterLibFunc(&sqliteInit, library, "sqlite3_abi_init")
	purego.RegisterLibFunc(&sqliteExec, library, "sqlite3_abi_exec")
	purego.RegisterLibFunc(&sqliteClose, library, "sqlite3_close")
	purego.RegisterLibFunc(&sqliteErrMsg, library, "sqlite3_errmsg")
	purego.RegisterLibFunc(&sqliteErrCode, library, "sqlite3_errcode")

	purego.RegisterLibFunc(&sqlitePrepareV2, library, "sqlite3_prepare_v2")
	purego.RegisterLibFunc(&sqliteFinalize, library, "sqlite3_finalize")
	purego.RegisterLibFunc(&sqliteReset, library, "sqlite3_reset")
	purego.RegisterLibFunc(&sqliteStep, library, "sqlite3_step")

	purego.RegisterLibFunc(&sqliteBindBlob, library, "sqlite3_bind_blob")
	purego.RegisterLibFunc(&sqliteBindBlob64, library, "sqlite3_bind_blob64")
	purego.RegisterLibFunc(&sqliteBindDouble, library, "sqlite3_bind_double")
	purego.RegisterLibFunc(&sqliteBindInt, library, "sqlite3_bind_int")
	purego.RegisterLibFunc(&sqliteBindInt64, library, "sqlite3_bind_int64")
	purego.RegisterLibFunc(&sqliteBindNull, library, "sqlite3_bind_null")
	purego.RegisterLibFunc(&sqliteBindText, library, "sqlite3_bind_text")

	purego.RegisterLibFunc(&sqlite3ColumnName, library, "sqlite3_column_name")
	purego.RegisterLibFunc(&sqlite3ColumnText, library, "sqlite3_column_text")

	return &CSQLiteImpl{
		libSQLite:         library,
		sqliteInit:        sqliteInit,
		sqliteExec:        sqliteExec,
		sqliteClose:       sqliteClose,
		sqliteErrMsg:      sqliteErrMsg,
		sqliteErrCode:     sqliteErrCode,
		sqlitePrepareV2:   sqlitePrepareV2,
		sqliteFinalize:    sqliteFinalize,
		sqliteBindBlob:    sqliteBindBlob,
		sqliteBindBlob64:  sqliteBindBlob64,
		sqliteBindDouble:  sqliteBindDouble,
		sqliteBindInt:     sqliteBindInt,
		sqliteBindInt64:   sqliteBindInt64,
		sqliteBindNull:    sqliteBindNull,
		sqliteBindText:    sqliteBindText,
		sqlite3ColumnName: sqlite3ColumnName,
		sqlite3ColumnText: sqlite3ColumnText,
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

func (c *CSQLiteImpl) SQLitePrepare(ctx *CSQLiteCtx, zSql string) *CSQLiteStmt {
	stmt := c.sqlitePrepareV2(ctx.ctx, zSql, len(zSql))
	return &CSQLiteStmt{stmt}
}

func (c *CSQLiteImpl) SQLiteFinalize(stmt *CSQLiteStmt) int {
	return c.sqliteFinalize(stmt.stmt)
}

func (c *CSQLiteImpl) SQLiteReset(stmt *CSQLiteStmt) int {
	return c.sqliteReset(stmt.stmt)
}

func (c *CSQLiteImpl) SQLiteStep(stmt *CSQLiteStmt) int {
	return c.sqliteStep(stmt.stmt)
}

func (c *CSQLiteImpl) SQLiteBindBlob(stmt *CSQLiteStmt, index int, zData []byte) int {
	return c.sqliteBindBlob(stmt.stmt, index, &zData[0], len(zData), 0)
}

func (c *CSQLiteImpl) SQLiteBindBlob64(stmt *CSQLiteStmt, index int, zData []byte) int {
	return c.sqliteBindBlob64(stmt.stmt, index, &zData[0], uint64(len(zData)), 0)
}

func (c *CSQLiteImpl) SQLiteBindDouble(stmt *CSQLiteStmt, index int, zData float64) int {
	return c.sqliteBindDouble(stmt.stmt, index, zData)
}

func (c *CSQLiteImpl) SQLiteBindInt(stmt *CSQLiteStmt, index int, zData int) int {
	return c.sqliteBindInt(stmt.stmt, index, zData)
}

func (c *CSQLiteImpl) SQLiteBindInt64(stmt *CSQLiteStmt, index int, zData int64) int {
	return c.sqliteBindInt64(stmt.stmt, index, zData)
}

func (c *CSQLiteImpl) SQLiteBindNull(stmt *CSQLiteStmt, index int) int {
	return c.sqliteBindNull(stmt.stmt, index)
}

func (c *CSQLiteImpl) SQLiteBindText(stmt *CSQLiteStmt, index int, zData string) int {
	return c.sqliteBindText(stmt.stmt, index, zData, len(zData))
}

func (c *CSQLiteImpl) SQLiteColumnName(stmt *CSQLiteStmt, index int) string {
	return c.sqlite3ColumnName(stmt.stmt, index)
}

func (c *CSQLiteImpl) SQLiteColumnText(stmt *CSQLiteStmt, index int) string {
	return c.sqlite3ColumnText(stmt.stmt, index)
}
