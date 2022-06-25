package repos

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"

	"euromoby.com/core/db"

	coremodels "euromoby.com/core/models"
)

const (
	connTimeout  = 1 * time.Second
	queryTimeout = 3 * time.Second
	txTimeout    = 1 * time.Second
)

type sqler interface {
	ToSql() (string, []interface{}, error)
}

func DBConnContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), connTimeout)
}

func DBTxContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), txTimeout)
}

func dbQueryContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), queryTimeout)
}

func Begin(conn db.Conn) (pgx.Tx, error) { //nolint:ireturn
	ctx, cancel := DBTxContext()
	defer cancel()

	return conn.Begin(ctx)
}

func Rollback(tx pgx.Tx) error {
	ctx, cancel := DBTxContext()
	defer cancel()

	return tx.Rollback(ctx)
}

func Commit(tx pgx.Tx) error {
	ctx, cancel := DBTxContext()
	defer cancel()

	return tx.Commit(ctx)
}

func DBQuerySingle(conn db.Conn, dst interface{}, sb sqler) error {
	stmt, args, err := sb.ToSql()
	if err != nil {
		return err
	}

	ctx, cancel := dbQueryContext()
	defer cancel()

	err = pgxscan.Get(ctx, conn, dst, stmt, args...)
	if pgxscan.NotFound(err) {
		return coremodels.ErrNotFound
	}

	return err
}

func DBQueryAll(conn db.Conn, dst interface{}, sb sqler) error {
	stmt, args, err := sb.ToSql()
	if err != nil {
		return err
	}

	ctx, cancel := dbQueryContext()
	defer cancel()

	return pgxscan.Select(ctx, conn, dst, stmt, args...)
}

func DBExec(conn db.Conn, sb sqler) error {
	stmt, args, err := sb.ToSql()
	if err != nil {
		return err
	}

	ctx, cancel := dbQueryContext()
	defer cancel()

	_, err = conn.Exec(ctx, stmt, args...)

	return err
}

func DBQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
