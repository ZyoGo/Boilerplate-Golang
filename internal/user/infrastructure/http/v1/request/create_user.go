package request

type CreateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
