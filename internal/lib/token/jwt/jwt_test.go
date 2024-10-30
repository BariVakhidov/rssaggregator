package jwtverifier

import (
	"reflect"
	"testing"
	"time"

	tokenverifier "github.com/BariVakhidov/rssaggregator/internal/lib/token"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWTVerifier_VerifyToken(t *testing.T) {
	type fields struct {
		secret     string
		email      string
		userId     uuid.UUID
		appId      uuid.UUID
		exp        time.Time
		signMethod jwt.SigningMethod
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name: "Happy path",
			fields: fields{
				secret:     gofakeit.LetterN(5),
				email:      gofakeit.Email(),
				userId:     uuid.New(),
				appId:      uuid.New(),
				exp:        time.Now().Add(time.Second * 5),
				signMethod: jwt.SigningMethodHS256,
			},
			wantErr: nil,
		},
		{
			name: "Expired token",
			fields: fields{
				secret:     gofakeit.LetterN(5),
				email:      gofakeit.Email(),
				userId:     uuid.New(),
				appId:      uuid.New(),
				exp:        time.Now().Add(-time.Second * 5),
				signMethod: jwt.SigningMethodHS256,
			},
			wantErr: tokenverifier.ErrInvalidToken,
		},
		{
			name: "Wrong signing type",
			fields: fields{
				secret:     gofakeit.LetterN(5),
				email:      gofakeit.Email(),
				userId:     uuid.New(),
				appId:      uuid.New(),
				exp:        time.Now().Add(time.Second * 5),
				signMethod: jwt.SigningMethodNone,
			},
			wantErr: tokenverifier.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fields := tt.fields

			j := New(fields.secret)

			payload := &tokenverifier.Payload{
				UserID: fields.userId,
				Email:  fields.email,
				AppID:  fields.appId,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(fields.exp),
				},
			}

			token, err := newToken(fields.signMethod, payload, fields.secret)
			if err != nil {
				t.Errorf("VerifyToken() create token error = %v", err)
				return
			}

			got, err := j.VerifyToken(token)
			if err != nil {
				if tt.wantErr != nil {
					assert.ErrorIs(t, err, tt.wantErr, "VerifyToken() error = %v, wantErr %v")
					return
				}
				t.Errorf("VerifyToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, payload) {
				t.Errorf("VerifyToken() got = %+v, want %+v", got, payload)
			}
		})
	}
}

func newToken(method jwt.SigningMethod, payload *tokenverifier.Payload, secret string) (string, error) {
	token := jwt.NewWithClaims(method, payload)

	var key any
	key = []byte(secret)
	if method == jwt.SigningMethodNone {
		key = jwt.UnsafeAllowNoneSignatureType
	}
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
