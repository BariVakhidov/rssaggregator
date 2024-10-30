package jwtverifier

import (
	tokenverifier "github.com/BariVakhidov/rssaggregator/internal/lib/token"
	"github.com/golang-jwt/jwt/v5"
)

type JWTVerifier struct {
	secret string
}

func New(secret string) *JWTVerifier {
	return &JWTVerifier{secret: secret}
}

func (j *JWTVerifier) VerifyToken(token string) (*tokenverifier.Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, tokenverifier.ErrInvalidToken
		}

		return []byte(j.secret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &tokenverifier.Payload{}, keyFunc)
	if err != nil {
		return nil, tokenverifier.ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*tokenverifier.Payload)
	if !ok {
		return nil, tokenverifier.ErrInvalidToken
	}

	return payload, nil
}
