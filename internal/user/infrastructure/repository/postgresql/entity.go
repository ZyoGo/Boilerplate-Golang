package postgresql

import (
	"github.com/ZyoGo/default-ddd-http/internal/user/core"
	"github.com/jackc/pgx/v5/pgtype"
)

type userEntity struct {
	ID       pgtype.Text
	Email    pgtype.Text
	Password pgtype.Text
}

func (ue *userEntity) toDomain() core.User {
	return core.User{
		ID:       ue.ID.String,
		Email:    ue.Email.String,
		Password: ue.Password.String,
	}
}
