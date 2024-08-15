package postgresql

import (
	"context"
	"errors"

	"github.com/ZyoGo/default-ddd-http/internal/user"
	"github.com/ZyoGo/default-ddd-http/pkg/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgreSQL struct {
	pool *pgxpool.Pool
}

type connCtx struct{}
type txCtx struct{}

type userEntity struct {
	ID       pgtype.Int8
	Email    pgtype.Text
	Password pgtype.Text
}

func (ue userEntity) toDomain() user.User {
	return user.User{
		ID:       ue.ID.Int64,
		Email:    ue.Email.String,
		Password: ue.Password.String,
	}
}

func NewPostgreSQL(pool *pgxpool.Pool) Repository {
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

func (repo *postgreSQL) InsertUser(ctx context.Context, user user.User) error {
	args := pgx.NamedArgs{
		"email":    user.Email,
		"password": user.Password,
	}

	queryInsert := `INSERT INTO users (email, password) VALUES (@email, @password)`

	// _, err := repo.db.Exec(ctx, queryInsert, args)
	_, err := repo.conn(ctx).Exec(ctx, queryInsert, args)
	if err != nil {
		return err
	}

	return nil
}

func (repo *postgreSQL) FindUserByEmail(ctx context.Context, email string) (user.User, error) {
	var ue userEntity
	args := pgx.NamedArgs{
		"email": email,
	}

	queryFind := `SELECT email FROM users WHERE email = @email`

	err := repo.conn(ctx).QueryRow(ctx, queryFind, args).Scan(&ue.ID, &ue.Email, &ue.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user.User{}, errors.New("data not found")
		}
		return user.User{}, err
	}

	return ue.toDomain(), nil
}

func (repo *postgreSQL) UpdateUser(ctx context.Context, user user.User) error {
	return errors.New("must be rollback")
	// args := pgx.NamedArgs{
	// 	"id":       user.ID,
	// 	"email":    user.Email,
	// 	"password": user.Password,
	// }

	// queryUpdate := `UPDATE users SET email = @email, password = @password WHERE id = @id`

	// _, err := repo.conn(ctx).Exec(ctx, queryUpdate, args)
	// if err != nil {
	// 	return err
	// }

	// return nil
}
