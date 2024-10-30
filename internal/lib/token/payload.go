package token

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	UserID uuid.UUID `json:"uid"`
	Email  string    `json:"email"`
	AppID  uuid.UUID `json:"app_id"`
	jwt.RegisteredClaims
}
