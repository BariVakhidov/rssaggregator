package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts an API key from
// the headers of an HTTP request
// Example:
// Authorization: ApiKey {insert apikey here}
func GetAPIKey(headers http.Header) (string, error) {
	key := headers.Get("Authorization")

	if key == "" {
		return key, errors.New("no authorization info found")
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
