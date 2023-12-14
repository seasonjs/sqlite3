package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/seasonjs/sqlite3"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	columnDate      string = "date"
	columnDatetime  string = "datetime"
	columnTimestamp string = "timestamp"
)

var driverName = "sqlite3"

func init() {
	sql.Register(driverName, &SQLiteDriver{})
}

// SQLiteDriver implements driver.Driver.
type SQLiteDriver struct {
	instance *sqlite3.SQLite3
}

// SQLiteConn implements driver.Conn and driver.ConnPrepareContext.
type SQLiteConn struct {
	db     *sqlite3.Connect
	mu     sync.Mutex
	loc    *time.Location
	txlock string
}

// SQLiteTx implements driver.Tx.
type SQLiteTx struct {
	conn *SQLiteConn
}

// SQLiteStmt implements driver.Stmt ,driver.StmtExecContext and driver.StmtQueryContext.
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

func (s *SQLiteDriver) Open(dsn string) (driver.Conn, error) {
	if s.instance == nil {
		sqlite, err := sqlite3.NewSQLite3Auto()
		if err != nil {
			return nil, err
		}
		s.instance = sqlite
	}
	var loc *time.Location
	txlock := "BEGIN"

	pos := strings.IndexRune(dsn, '?')
	if pos >= 1 {
		params, err := url.ParseQuery(dsn[pos+1:])
		if err != nil {
			return nil, err
		}
		if val := params.Get("_loc"); val != "" {
			switch strings.ToLower(val) {
			case "auto":
				loc = time.Local
			default:
				loc, err = time.LoadLocation(val)
				if err != nil {
					return nil, fmt.Errorf("invalid _loc: %v: %v", val, err)
				}
			}
		}
		// _mutex
		//if val := params.Get("_mutex"); val != "" {
		//	switch strings.ToLower(val) {
		//	case "no":
		//		mutex = C.SQLITE_OPEN_NOMUTEX
		//	case "full":
		//		mutex = C.SQLITE_OPEN_FULLMUTEX
		//	default:
		//		return nil, fmt.Errorf("invalid _mutex: %v", val)
		//	}
		//}

		// _txlock
		if val := params.Get("_txlock"); val != "" {
			switch strings.ToLower(val) {
			case "immediate":
				txlock = "BEGIN IMMEDIATE"
			case "exclusive":
				txlock = "BEGIN EXCLUSIVE"
			case "deferred":
				txlock = "BEGIN"
			default:
				return nil, fmt.Errorf("invalid _txlock: %v", val)
			}
		}

		if !strings.HasPrefix(dsn, "file:") {
			dsn = dsn[:pos]
		}

	}
	db, err := s.instance.Open(dsn)
	if err != nil {
		return nil, err
	}

	conn := &SQLiteConn{
		db:     db,
		loc:    loc,
		txlock: txlock,
	}

	runtime.SetFinalizer(conn, (*SQLiteConn).Close)
	return conn, nil
}

// Prepare returns a prepared statement, bound to this connection.
func (s *SQLiteConn) Prepare(query string) (driver.Stmt, error) {
	panic("implement me")
}

func (s *SQLiteConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SQLiteConn) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.db = nil
	s.mu.Unlock()
	runtime.SetFinalizer(s, nil)
	return nil
}

func (s *SQLiteConn) Begin() (driver.Tx, error) {
	if _, err := s.db.Exec(s.txlock); err != nil {
		return nil, err
	}
	return &SQLiteTx{conn: s}, nil
}

func (s *SQLiteTx) Commit() error {
	_, err := s.conn.db.Exec("COMMIT")
	if err != nil {
		_, err = s.conn.db.Exec("ROLLBACK")
	}
	return err
}

func (s *SQLiteTx) Rollback() error {
	_, err := s.conn.db.Exec("ROLLBACK")
	return err
}

func (s *SQLiteStmt) Close() error {
	//TODO implement me
	return nil
}

func (s *SQLiteStmt) NumInput() int {
	//TODO implement me
	return -1
}

func (s *SQLiteStmt) Exec(args []driver.Value) (driver.Result, error) {
	list := make([]driver.NamedValue, len(args))
	for i, v := range args {
		list[i] = driver.NamedValue{
			Ordinal: i + 1,
			Value:   v,
		}
	}
	return s.ExecContext(context.Background(), list)
}

func (s *SQLiteStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SQLiteStmt) Query(args []driver.Value) (driver.Rows, error) {
	list := make([]driver.NamedValue, len(args))
	for i, v := range args {
		list[i] = driver.NamedValue{
			Ordinal: i + 1,
			Value:   v,
		}
	}
	return s.QueryContext(context.Background(), list)
}

func (s *SQLiteStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SQLiteResult) LastInsertId() (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SQLiteResult) RowsAffected() (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SQLiteRows) Columns() []string {
	//TODO implement me
	panic("implement me")
}

func (s *SQLiteRows) Close() error {
	//TODO implement me
	panic("implement me")
}

func (s *SQLiteRows) Next(dest []driver.Value) error {
	//TODO implement me
	panic("implement me")
}
