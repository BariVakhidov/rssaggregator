package storage

import "errors"

const (
	UniqueViolationCode       = "23505"
	ForeignKeyViolationCode   = "23503"
	InvalidTextRepresentation = "22P02"
)

var (
	ErrUserExists          = errors.New("user exists")
	ErrFeedExists          = errors.New("feed exists")
	ErrFeedFollowExists    = errors.New("feed follow exists")
	ErrPostExists          = errors.New("post exists")
	ErrPendingUserExists   = errors.New("pending user exists")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrInvalidParams       = errors.New("invalid params")
	ErrNotFound            = errors.New("no rows")
)
