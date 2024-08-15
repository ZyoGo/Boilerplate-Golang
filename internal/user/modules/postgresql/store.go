package postgresql

import (
	"context"

	"github.com/ZyoGo/default-ddd-http/internal/user"
)

type Repository interface {
	// TransactionContext returns a copy of the parent context which begins a transaction
	// to PostgreSQL.
	TransactionContext(ctx context.Context) (context.Context, error)
	// Commit transaction from context.
	Commit(ctx context.Context) error
	// Rollback transaction from context.
	Rollback(ctx context.Context) error

	InsertUser(ctx context.Context, user user.User) error
	FindUserByEmail(ctx context.Context, email string) (user.User, error)
	UpdateUser(ctx context.Context, user user.User) error
}
