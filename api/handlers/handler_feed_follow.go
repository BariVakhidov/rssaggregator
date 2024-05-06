package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BariVakhidov/rssaggregator/db"

	"github.com/BariVakhidov/rssaggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) HandlerFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := &parameters{}

	err := decoder.Decode(params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feedFollow, dbErr := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    params.FeedID,
	})

	if dbErr != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error while creating feed follow: %v", dbErr))
		return
	}

	respondWithJson(w, http.StatusCreated, db.DatabaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *ApiConfig) HandlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, dbErr := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)

	if dbErr != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error while getting feed follows: %v", dbErr))
		return
	}

	respondWithJson(w, http.StatusCreated, db.DatabaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *ApiConfig) HandlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to parse feedFollowID: %v", err))
		return
	}

	dbErr := apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		ID:     feedFollowID,
	})

	if dbErr != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error while deleting feed follow: %v", dbErr))
		return
	}

	respondWithJson(w, http.StatusNoContent, struct{}{})
}
