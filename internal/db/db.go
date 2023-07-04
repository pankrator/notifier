package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"

	"github.com/pankrator/notifier/internal/config"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

const (
	_pgDriver   = "postgres"
	_connString = "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"
)

type DBConn struct {
	*sqlx.DB
}

type Logger struct {
	*log.Logger
}

func NewDBConn(ctx context.Context, c config.DB) (*DBConn, error) {
	dsn := fmt.Sprintf(_connString, c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)

	db, err := sql.Open(_pgDriver, dsn)
	if err != nil {
		return nil, err
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	zlogger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	loggerOptions := []sqldblogger.Option{
		sqldblogger.WithSQLQueryFieldname("sql"),
		sqldblogger.WithWrapResult(false),
		sqldblogger.WithExecerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithQueryerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithPreparerLevel(sqldblogger.LevelDebug),
	}

	db = sqldblogger.OpenDriver(dsn, db.Driver(), zerologadapter.New(zlogger), loggerOptions...)

	sqlxDB := sqlx.NewDb(db, _pgDriver)
	err = sqlxDB.Ping()
	if err != nil {
		return nil, err
	}

	return &DBConn{
		DB: sqlxDB,
	}, nil
}

func (t *DBConn) RunInTx(ctx context.Context, run func(ctx context.Context) error) error {
	tx, err := t.beginTx(ctx)
	if err != nil {
		return err
	}

	var done bool

	defer func() {
		if !done {
			_ = tx.Rollback()
		}
	}()

	if err := run(putTransactionToContext(ctx, tx)); err != nil {
		return err
	}

	done = true

	return tx.Commit()
}

func (t *DBConn) beginTx(ctx context.Context) (*Tx, error) {
	if tx := ExtractTransactionFromContext(ctx); tx != nil {
		return tx, nil
	}

	tx, err := t.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, err
	}

	txx := &Tx{
		Tx: tx,
	}
	return txx, nil
}

func (conn *DBConn) QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	tx := ExtractTransactionFromContext(ctx)
	if tx != nil {
		return tx.QueryxContext(ctx, query, args...)
	}

	return conn.DB.QueryxContext(ctx, query, args...)
}

func (conn *DBConn) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	tx := ExtractTransactionFromContext(ctx)
	if tx != nil {
		return tx.NamedExecContext(ctx, query, arg)
	}

	return conn.DB.NamedExecContext(ctx, query, arg)
}

func (conn *DBConn) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx := ExtractTransactionFromContext(ctx)
	if tx != nil {
		return tx.SelectContext(ctx, dest, query, args...)
	}

	return conn.DB.SelectContext(ctx, dest, query, args...)
}

func (conn *DBConn) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	tx := ExtractTransactionFromContext(ctx)
	if tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}

	return conn.DB.ExecContext(ctx, query, args...)
}
