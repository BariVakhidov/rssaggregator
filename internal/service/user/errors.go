package userservice

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidUserID      = errors.New("invalid userID")
	ErrUserExists         = errors.New("user exists")
	ErrUserLocked         = errors.New("user temporary locked")
	ErrUserNotFound       = errors.New("user not found")
)
