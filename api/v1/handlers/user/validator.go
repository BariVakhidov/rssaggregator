package user

import "github.com/BariVakhidov/rssaggregator/internal/model"

func (i *Implementation) validateUserInfo(userInfo model.UserInfo) map[string]string {
	errors := make(map[string]string)

	if len(userInfo.Email) == 0 {
		errors["email"] = ErrEmailRequired.Error()
	}

	if err := i.validator.Var(userInfo.Email, "email"); err != nil {
		errors["email"] = ErrEmailInvalid.Error()
	}

	if len(userInfo.Password) == 0 {
		errors["password"] = ErrPasswordRequired.Error()
	}

	return errors
}
