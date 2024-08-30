package core

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")

	// Error for password
	ErrPasswordTooShort         = errors.New("password is too short; it must be at least 6 characters long")
	ErrPasswordMissingUppercase = errors.New("password must contain at least one uppercase letter")
	ErrPasswordMissingLowercase = errors.New("password must contain at least one lowercase letter")
	ErrPasswordMissingDigit     = errors.New("password must contain at least one numeric digit")
	ErrPasswordMissingSpecial   = errors.New("password must contain at least one special character")
)
