package postgresql

import (
	"github.com/ZyoGo/default-ddd-http/internal/user"
	"github.com/jackc/pgx/v5/pgtype"
)

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
