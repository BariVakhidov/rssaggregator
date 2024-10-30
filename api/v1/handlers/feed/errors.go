package feed

import "errors"

var (
	ErrFeedExists         = errors.New("feed already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrFeedNotFound       = errors.New("feed not found")
	ErrFeedNameRequired   = errors.New("feed name required")
	ErrFeedUrlRequired    = errors.New("feed url required")
)
