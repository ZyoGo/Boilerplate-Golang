package core

import "context"

type Repository interface {
	// TransactionContext returns a copy of the parent context which begins a transaction
	// to PostgreSQL.
	TransactionContext(ctx context.Context) (context.Context, error)
	// Commit transaction from context.
	Commit(ctx context.Context) error
	// Rollback transaction from context.
	Rollback(ctx context.Context) error

	InsertUser(ctx context.Context, user User) error
	FindUserByEmail(ctx context.Context, email string) (User, error)
	UpdateUser(ctx context.Context, user User) error
}
