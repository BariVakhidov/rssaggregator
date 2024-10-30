package feed

import (
	"encoding/json"
	"net/http"

	handlersv1 "github.com/BariVakhidov/rssaggregator/api/v1/handlers"
	auth "github.com/BariVakhidov/rssaggregator/internal/delivery/http/middleware"
	jsonlib "github.com/BariVakhidov/rssaggregator/internal/lib/json"
	"github.com/BariVakhidov/rssaggregator/internal/model"
)

// @Summary Create a new feed
// @Description Add a new RSS feed to the system.
// @Tags Feeds
// @Accept json
// @Produce json
// @Param feed body model.FeedInfo true "Feed Request Body"
// @Success 200 {object} model.Feed
// @Failure 400 {object} model.APIErr "Invalid input or validation errors"
// @Failure 500 {object} model.APIErr "Internal server error"
// @Router /feeds [post]
func (i *Implementation) CreateFeed(w http.ResponseWriter, r *http.Request) {
	var feedInfo model.FeedInfo

	if err := json.NewDecoder(r.Body).Decode(&feedInfo); err != nil {
		jsonlib.RespondWithError(w, i.log, model.JSONErr())
		return
	}

	if validateErrors := i.validateFeedInfo(feedInfo); len(validateErrors) > 0 {
		jsonlib.RespondWithError(w, i.log, model.InvalidRequestData(validateErrors))
		return
	}

	createdFeed, err := i.feedService.CreateFeed(r.Context(), feedInfo)
	if err != nil {
		jsonlib.RespondWithError(w, i.log, i.checkFeedServiceErr(err))
		return
	}

	jsonlib.RespondWithJson(w, i.log, http.StatusOK, createdFeed)
}

// @Summary Create a new feed follow
// @Description Add a new RSS feed follow to the user.
// @Tags Feeds
// @Accept json
// @Produce json
// @Param feed body model.FeedFollowInfo true "Feed Request Body"
// @Success 200 {object} model.FeedFollow
// @Failure 400 {object} model.APIErr "Invalid input or validation errors"
// @Failure 500 {object} model.APIErr "Internal server error"
// @Router /feed_follows [post]
func (i *Implementation) CreateFeedFollow(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		jsonlib.RespondWithError(w, i.log, model.NewAPIErr(http.StatusUnauthorized, handlersv1.ErrUnauthorized))
		return
	}

	var followInfo model.FeedFollowInfo

	if err := json.NewDecoder(r.Body).Decode(&followInfo); err != nil {
		jsonlib.RespondWithError(w, i.log, model.JSONErr())
		return
	}

	follow, err := i.feedService.CreateFeedFollow(r.Context(), userID, followInfo.FeedId)
	if err != nil {
		jsonlib.RespondWithError(w, i.log, i.checkFeedServiceErr(err))
		return
	}

	jsonlib.RespondWithJson(w, i.log, http.StatusOK, follow)
}
