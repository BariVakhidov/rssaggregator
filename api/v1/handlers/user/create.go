package user

import (
	"encoding/json"
	"net/http"

	jsonlib "github.com/BariVakhidov/rssaggregator/internal/lib/json"
	"github.com/BariVakhidov/rssaggregator/internal/model"
)

// @Summary Register user
// @Description Create new user with credentials.
// @Tags User
// @Accept json
// @Produce json
// @Param user body model.UserInfo true "User Request Body"
// @Success 200 {object} nil
// @Failure 404 {object} model.APIErr
// @Router /users [post]
func (i *Implementation) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userInfo model.UserInfo

	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		jsonlib.RespondWithError(w, i.log, model.JSONErr())
		return
	}

	if errors := i.validateUserInfo(userInfo); len(errors) > 0 {
		jsonlib.RespondWithError(w, i.log, model.InvalidRequestData(errors))
	}

	if err := i.userService.CreateUser(r.Context(), userInfo); err != nil {
		jsonlib.RespondWithError(w, i.log, i.userServiceErr(err))
		return
	}

	jsonlib.RespondWithJson(w, i.log, http.StatusCreated, nil)
}
