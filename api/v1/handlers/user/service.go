package user

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	handlersv1 "github.com/BariVakhidov/rssaggregator/api/v1/handlers"

	"github.com/BariVakhidov/rssaggregator/internal/model"
	userservice "github.com/BariVakhidov/rssaggregator/internal/service/user"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, userInfo model.UserInfo) error
	Login(ctx context.Context, userInfo model.UserInfo) (string, error)
	User(ctx context.Context, userID uuid.UUID) (*model.User, error)
}

type Implementation struct {
	userService UserService
	log         *slog.Logger
	validator   *validator.Validate
}

func New(log *slog.Logger, userService UserService) *Implementation {
	return &Implementation{
		userService: userService,
		log:         log,
		validator:   validator.New(),
	}
}

func (i *Implementation) userServiceErr(err error) error {
	if errors.Is(err, userservice.ErrInvalidCredentials) {
		return model.NewAPIErr(http.StatusBadRequest, handlersv1.ErrInvalidCredentials)
	}

	if errors.Is(err, userservice.ErrUserLocked) {
		return model.NewAPIErr(http.StatusLocked, ErrUserLocked)
	}

	if errors.Is(err, userservice.ErrUserExists) {
		return model.NewAPIErr(http.StatusConflict, ErrUserExists)
	}

	return err
}
