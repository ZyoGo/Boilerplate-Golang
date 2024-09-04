package bcrypt

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type service struct{}

func New() *service {
	return &service{}
}

func (b service) HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordHash), nil
}

func (b service) IsSamePassword(passwordStr, passwordHash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordStr), []byte(passwordHash))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
