package user

import "errors"

var (
	ErrEmailRequired    = errors.New("email is required")
	ErrEmailInvalid     = errors.New("email is invalid")
	ErrPasswordRequired = errors.New("password required")
	ErrNameRequired     = errors.New("name is required")
	ErrUserLocked       = errors.New("user temporary locked")
	ErrUserExists       = errors.New("user already exists")
)
