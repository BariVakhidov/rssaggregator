package user

import (
	"encoding/json"
	"net/http"

	jsonlib "github.com/BariVakhidov/rssaggregator/internal/lib/json"
	"github.com/BariVakhidov/rssaggregator/internal/model"
)

// @Summary Login user
// @Description Login user.
// @Tags User
// @Accept json
// @Produce json
// @Param user body model.UserInfo true "User Request Body"
// @Success 200 {object} string
// @Failure 401 {object} model.APIErr
// @Failure 400 {object} model.APIErr
// @Failure 423 {object} model.APIErr
// @Router /uses/login [post]
func (i *Implementation) Login(w http.ResponseWriter, r *http.Request) {
	var userInfo model.UserInfo

	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		jsonlib.RespondWithError(w, i.log, model.JSONErr())
		return
	}

	if errors := i.validateUserInfo(userInfo); len(errors) > 0 {
		jsonlib.RespondWithError(w, i.log, model.InvalidRequestData(errors))
	}

	token, err := i.userService.Login(r.Context(), userInfo)
	if err != nil {
		jsonlib.RespondWithError(w, i.log, i.userServiceErr(err))
		return
	}
	//TODO: add cookies for auth

	jsonlib.RespondWithJson(w, i.log, http.StatusCreated, token)
}
