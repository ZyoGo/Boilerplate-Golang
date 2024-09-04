package core

type Hash interface {
	HashPassword(password string) (string, error)
	IsSamePassword(passwordStr, passwordHash string) (bool, error)
}
