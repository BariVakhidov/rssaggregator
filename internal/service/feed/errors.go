package feed

import "errors"

var (
	ErrFeedExists         = errors.New("feed exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound           = errors.New("not found")
)
