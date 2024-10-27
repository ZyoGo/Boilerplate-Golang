package core

type User struct {
	ID       string
	Email    string
	Password string
}

type Auth struct {
	ID          string
	Email       string
	AccessToken string
}
