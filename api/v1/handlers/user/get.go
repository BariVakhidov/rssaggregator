package user

import (
	"net/http"

	handlersv1 "github.com/BariVakhidov/rssaggregator/api/v1/handlers"
	auth "github.com/BariVakhidov/rssaggregator/internal/delivery/http/middleware"
	jsonlib "github.com/BariVakhidov/rssaggregator/internal/lib/json"
	"github.com/BariVakhidov/rssaggregator/internal/model"
)

// @Summary Get user by ID from auth_token
// @Description Retrieve user information by user ID.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.User
// @Failure 404 {object} model.APIErr
// @Router /users [get]
func (i *Implementation) User(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		jsonlib.RespondWithError(w, i.log, model.NewAPIErr(http.StatusUnauthorized, handlersv1.ErrUnauthorized))
		return
	}

	user, err := i.userService.User(r.Context(), userID)
	if err != nil {
		jsonlib.RespondWithError(w, i.log, i.userServiceErr(err))
		return
	}

	jsonlib.RespondWithJson(w, i.log, http.StatusOK, user)
}
