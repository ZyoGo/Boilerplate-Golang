package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/ZyoGo/default-ddd-http/internal/user/core"
	"github.com/ZyoGo/default-ddd-http/pkg/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type connCtx struct{}
type txCtx struct{}

type postgreSQL struct {
	pool *pgxpool.Pool
}

func NewPostgreSQL(pool *pgxpool.Pool) core.Repository {
	return &postgreSQL{pool: pool}
}

func (db *postgreSQL) TransactionContext(ctx context.Context) (context.Context, error) {
	tx, err := db.conn(ctx).Begin(ctx)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, txCtx{}, tx), nil
}

func (db *postgreSQL) Commit(ctx context.Context) error {
	if tx, ok := ctx.Value(txCtx{}).(pgx.Tx); ok && tx != nil {
		return tx.Commit(ctx)
	}
	return errors.New("context has no transaction")
}

func (db *postgreSQL) Rollback(ctx context.Context) error {
	if tx, ok := ctx.Value(txCtx{}).(pgx.Tx); ok && tx != nil {
		return tx.Rollback(ctx)
	}
	return errors.New("context has no transaction")
}

// conn returns a PostgreSQL transaction if one exists.
// If not, returns a connection if a connection has been acquired by calling WithAcquire.
// Otherwise, it returns *pgxpool.Pool which acquires the connection and closes it immediately after a SQL command is executed.
func (db *postgreSQL) conn(ctx context.Context) database.PGXQuerier {
	if tx, ok := ctx.Value(txCtx{}).(pgx.Tx); ok && tx != nil {
		return tx
	}
	if res, ok := ctx.Value(connCtx{}).(*pgxpool.Conn); ok && res != nil {
		return res
	}
	return db.pool
}

func (repo *postgreSQL) InsertUser(ctx context.Context, user core.User) error {
	args := pgx.NamedArgs{
		"email":    user.Email,
		"password": user.Password,
	}

	_, err := repo.conn(ctx).Exec(ctx, queryInsertUser, args)
	if err != nil {
		return err
	}

	return nil
}

func (repo *postgreSQL) FindUserByEmail(ctx context.Context, email string) (core.User, error) {
	var ue userEntity
	args := pgx.NamedArgs{
		"email": email,
	}

	err := repo.conn(ctx).QueryRow(ctx, queryFindUserByEmail, args).Scan(&ue.Email)
	if err != nil && err != pgx.ErrNoRows {
		return core.User{}, err
	}

	return ue.toDomain(), nil
}

func (repo *postgreSQL) UpdateUser(ctx context.Context, user core.User) error {
	args := pgx.NamedArgs{
		"id":       user.ID,
		"password": user.Password,
	}

	_, err := repo.conn(ctx).Exec(ctx, queryUpdateUser, args)
	if err != nil {
		return err
	}

	return nil
}

func (repo *postgreSQL) FindUsers(ctx context.Context, filter core.FindUserFilter) ([]core.User, error) {
	args := pgx.NamedArgs{}

	if filter.Email != "" {
		args["email"] = filter.Email
	}

	fmt.Println("args = ", args)

	return nil, nil
}

func (repo *postgreSQL) FindUserByID(ctx context.Context, id string) (core.User, error) {
	return core.User{}, nil
}
