package feed

import (
	"net/http"
	"strconv"

	handlersv1 "github.com/BariVakhidov/rssaggregator/api/v1/handlers"
	auth "github.com/BariVakhidov/rssaggregator/internal/delivery/http/middleware"
	jsonlib "github.com/BariVakhidov/rssaggregator/internal/lib/json"
	"github.com/BariVakhidov/rssaggregator/internal/model"
)

// Feeds retrieves all available feeds.
// @Summary Get all feeds
// @Description Retrieve a list of all feeds.
// @Tags Feeds
// @Produce json
// @Success 200 {array} model.Feed "List of feeds"
// @Failure 500 {object} model.APIErr "Internal server error"
// @Router /feeds [get]
func (i *Implementation) Feeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := i.feedService.Feeds(r.Context())
	if err != nil {
		jsonlib.RespondWithError(w, i.log, i.checkFeedServiceErr(err))
		return
	}

	jsonlib.RespondWithJson(w, i.log, http.StatusOK, feeds)
}

// FeedFollows retrieves feed follows for the authenticated user.
// @Summary Get followed feeds
// @Description Retrieve a list of feeds followed by the user.
// @Tags Feeds
// @Produce json
// @Param limit query int false "Number of feeds to return (default: 10)"
// @Success 200 {array} model.FeedFollow "List of followed feeds"
// @Failure 400 {object} model.APIErr "Invalid request or parameters"
// @Failure 401 {object} model.APIErr "Unauthorized"
// @Failure 500 {object} model.APIErr "Internal server error"
// @Router /feed_follows [get]

func (i *Implementation) FeedFollows(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		jsonlib.RespondWithError(w, i.log, model.NewAPIErr(http.StatusUnauthorized, handlersv1.ErrUnauthorized))
		return
	}

	var limit int
	limitString := r.URL.Query().Get("limit")

	if len(r.URL.Query()) != 0 && len(limitString) == 0 {
		jsonlib.RespondWithError(w, i.log, model.NewAPIErr(http.StatusBadRequest, ErrInvalidCredentials))
		return
	}

	if len(limitString) == 0 {
		limit = 10
	} else {
		limit, err := strconv.Atoi(limitString)
		if err != nil || limit == 0 {
			jsonlib.RespondWithError(w, i.log, model.NewAPIErr(http.StatusBadRequest, ErrInvalidCredentials))
			return
		}
	}

	follows, err := i.feedService.FeedFollows(r.Context(), model.FeedFollowsInfo{
		UserID: userID,
		Limit:  limit,
	})
	if err != nil {
		jsonlib.RespondWithError(w, i.log, i.checkFeedServiceErr(err))
		return
	}

	jsonlib.RespondWithJson(w, i.log, http.StatusOK, follows)
}
