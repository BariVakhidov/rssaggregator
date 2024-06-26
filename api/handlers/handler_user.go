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

func (apiCfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	newUser, dbErr := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if dbErr != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error while creating user: %v", dbErr))
		return
	}

	respondWithJson(w, http.StatusCreated, db.DatabaseUserToUser(newUser))
}

func (apiCfg *ApiConfig) HandlerGetUserByAPIKey(w http.ResponseWriter, _ *http.Request, user database.User) {
	respondWithJson(w, http.StatusOK, db.DatabaseUserToUser(user))
}

func (apiCfg *ApiConfig) HandlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostsForUser(
		r.Context(),
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  100,
		},
	)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't get posts: %v", err))
		return
	}

	respondWithJson(w, http.StatusOK, db.DatabasePostsToPosts(posts))
}
