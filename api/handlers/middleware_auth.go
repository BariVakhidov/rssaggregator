package handlers

import (
	"fmt"
	"net/http"

	"github.com/BariVakhidov/rssaggregator/internal/auth"
	"github.com/BariVakhidov/rssaggregator/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *ApiConfig) MiddlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil {
			respondWithError(w, http.StatusForbidden, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, dbErr := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)

		if dbErr != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error getting user: %v", dbErr))
			return
		}

		handler(w, r, user)
	}
}
