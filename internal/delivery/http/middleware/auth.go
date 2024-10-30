package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/BariVakhidov/rssaggregator/internal/lib/logger/sl"
	tokenverifier "github.com/BariVakhidov/rssaggregator/internal/lib/token"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type ContextKey string

const (
	ErrorKey  ContextKey = "error"
	UserIDKey ContextKey = "userID"
)

// extractToken extracts an API key from
// the headers of an HTTP request
// Example:
// Authorization: ApiKey {insert apikey here}
func extractToken(headers http.Header) (string, error) {
	key := headers.Get("Authorization")

	if key == "" {
		return key, nil
	}

	values := strings.Split(key, " ")

	if len(values) != 2 {
		return "", errors.New("malformed auth header")
	}

	if values[0] != "ApiKey" {
		return "", errors.New("malformed auth header")
	}

	return values[1], nil
}

func New(log *slog.Logger, tokenVerifier tokenverifier.Verifier) func(next http.Handler) http.Handler {
	const op = "middleware.auth.New"
	log = log.With(slog.String("op", op))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := extractToken(r.Header)
			if err != nil {
				log.Warn("failed to parse token from header", sl.Err(err))
				ctx := context.WithValue(r.Context(), ErrorKey, ErrInvalidToken)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if len(token) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			payload, err := tokenVerifier.VerifyToken(token)
			if err != nil {
				log.Warn("failed to verify token", sl.Err(err))
				ctx := context.WithValue(r.Context(), ErrorKey, ErrInvalidToken)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, payload.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

func ErrorFromContext(ctx context.Context) (error, bool) {
	err, ok := ctx.Value(ErrorKey).(error)
	return err, ok
}
