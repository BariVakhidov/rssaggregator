package model

import (
	"errors"
	"fmt"
	"net/http"
)

type APIErr struct {
	StatusCode int `json:"status_code"`
	Msg        any `json:"msg"`
}

func (a APIErr) Error() string {
	return fmt.Sprintf("api error %d", a.StatusCode)
}

func NewAPIErr(code int, err error) APIErr {
	return APIErr{
		StatusCode: code,
		Msg:        err.Error(),
	}
}

var (
	ErrInvalidJSON = errors.New("invalid json")
)

func JSONErr() APIErr {
	return NewAPIErr(http.StatusBadRequest, ErrInvalidJSON)
}

func InvalidRequestData(errors map[string]string) APIErr {
	return APIErr{
		StatusCode: http.StatusBadRequest,
		Msg:        errors,
	}
}
