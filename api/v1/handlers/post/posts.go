package post

import (
	"errors"
	"net/http"
	"strconv"

	handlersv1 "github.com/BariVakhidov/rssaggregator/api/v1/handlers"
	auth "github.com/BariVakhidov/rssaggregator/internal/delivery/http/middleware"
	jsonlib "github.com/BariVakhidov/rssaggregator/internal/lib/json"
	"github.com/BariVakhidov/rssaggregator/internal/model"
	postService "github.com/BariVakhidov/rssaggregator/internal/service/post"
)

// @Summary Get user posts
// @Description Retrieve user posts.
// @Tags Post
// @Accept json
// @Produce json
// @Success 200 {array} model.Post  "List of posts of followed feeds"
// @Failure 404 {object} model.APIErr
// @Router /posts [get]
func (i *Implementation) Posts(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		jsonlib.RespondWithError(w, i.log, model.NewAPIErr(http.StatusUnauthorized, handlersv1.ErrUnauthorized))
		return
	}

	var limit int
	limitString := r.URL.Query().Get("limit")

	if len(r.URL.Query()) != 0 && len(limitString) == 0 {
		jsonlib.RespondWithError(w, i.log, model.NewAPIErr(http.StatusBadRequest, handlersv1.ErrInvalidCredentials))
		return
	}

	if len(limitString) == 0 {
		limit = 10
	} else {
		limit, err := strconv.Atoi(limitString)
		if err != nil || limit == 0 {
			jsonlib.RespondWithError(w, i.log, model.NewAPIErr(http.StatusBadRequest, handlersv1.ErrInvalidCredentials))
			return
		}
	}

	posts, err := i.postService.Posts(r.Context(), model.PostInfo{UserId: userID, Limit: limit})
	if err != nil {
		returnedErr := err
		if errors.Is(err, postService.ErrInvalidCredentials) {
			returnedErr = model.NewAPIErr(http.StatusBadRequest, handlersv1.ErrInvalidCredentials)
		}
		jsonlib.RespondWithError(w, i.log, returnedErr)
		return
	}

	jsonlib.RespondWithJson(w, i.log, http.StatusOK, posts)
}
