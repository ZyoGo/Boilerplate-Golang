package postgresql

import (
	"github.com/ZyoGo/default-ddd-http/internal/user-v1/core"
	"github.com/jackc/pgx/v5/pgtype"
)

type userEntity struct {
	ID       pgtype.Int8
	Email    pgtype.Text
	Password pgtype.Text
}

func (ue userEntity) toDomain() core.User {
	return core.User{
		ID:       ue.ID.Int64,
		Email:    ue.Email.String,
		Password: ue.Password.String,
	}
}
