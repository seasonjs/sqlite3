package sqlite3

type SQLiteDriver struct {
}

// SQLiteConn implements driver.Conn.
type SQLiteConn struct {
	//mu sync.Mutex
}

// SQLiteTx implements driver.Tx.
type SQLiteTx struct {
	//c *SQLiteConn
}

// SQLiteStmt implements driver.Stmt.
type SQLiteStmt struct {
	//mu sync.Mutex
	//c  *SQLiteConn
}

// SQLiteResult implements sql.Result.
type SQLiteResult struct {
	//id      int64
	//changes int64
}

// SQLiteRows implements driver.Rows.
type SQLiteRows struct {
	//s        *SQLiteStmt
	//nc       int
	//cols     []string
	//decltype []string
	//cls      bool
	//closed   bool
	//ctx      context.Context // no better alternative to pass context into Next() method
}
