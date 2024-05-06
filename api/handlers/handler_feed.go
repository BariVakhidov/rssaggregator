package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BariVakhidov/rssaggregator/db"

	"github.com/BariVakhidov/rssaggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) HandlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)

	params := &parameters{}

	err := decoder.Decode(params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	newFeed, dbErr := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url:       params.Url,
	})

	if dbErr != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error while creating feed: %v", dbErr))
		return
	}

	respondWithJson(w, http.StatusCreated, db.DatabaseFeedToFeed(newFeed))
}

func (apiCfg *ApiConfig) HandlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, dbErr := apiCfg.DB.GetFeeds(r.Context())

	if dbErr != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error while getting feeds: %v", dbErr))
		return
	}

	respondWithJson(w, http.StatusCreated, db.DatabaseFeedsToFeeds(feeds))
}
