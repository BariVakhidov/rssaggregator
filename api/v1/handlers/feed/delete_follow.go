package feed

import (
	"net/http"

	handlersv1 "github.com/BariVakhidov/rssaggregator/api/v1/handlers"
	auth "github.com/BariVakhidov/rssaggregator/internal/delivery/http/middleware"
	jsonlib "github.com/BariVakhidov/rssaggregator/internal/lib/json"
	"github.com/BariVakhidov/rssaggregator/internal/model"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// @Summary Delete a feed follow
// @Description Unfollow a feed for the authenticated user.
// @Tags Feeds
// @Param feedFollowID path string true "Feed Follow ID"
// @Produce json
// @Success 200 {object} nil "Successfully unfollowed the feed"
// @Failure 400 {object} model.APIErr "Invalid request"
// @Failure 401 {object} model.APIErr "Unauthorized"
// @Failure 500 {object} model.APIErr "Internal server error"
// @Router /feed_follows/{feedFollowID} [delete]
func (i *Implementation) DeleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		jsonlib.RespondWithError(w, i.log, model.NewAPIErr(http.StatusUnauthorized, handlersv1.ErrUnauthorized))
		return
	}

	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	if len(feedFollowIDStr) == 0 {
		jsonlib.RespondWithError(w, i.log, model.NewAPIErr(http.StatusBadRequest, ErrInvalidCredentials))
		return
	}

	followId, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		jsonlib.RespondWithError(w, i.log, model.NewAPIErr(http.StatusBadRequest, ErrInvalidCredentials))
		return
	}

	if err := i.feedService.DeleteFeedFollow(r.Context(), userID, followId); err != nil {
		jsonlib.RespondWithError(w, i.log, i.checkFeedServiceErr(err))
		return
	}

	jsonlib.RespondWithJson(w, i.log, http.StatusOK, nil)
}
