package userservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/BariVakhidov/rssaggregator/internal/lib/logger/sl"
	"github.com/BariVakhidov/rssaggregator/internal/model"
	"github.com/BariVakhidov/rssaggregator/internal/storage"
	ssov1 "github.com/BariVakhidov/ssoprotos/gen/go/sso"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserStorage interface {
	SavePendingUser(ctx context.Context, pendingUserInfo model.PendingUserInfo) (model.PendingUser, error)
	PendingUserByEmail(ctx context.Context, userEmail string) (model.PendingUser, error)
	ChangePendingUserStatus(ctx context.Context, userId uuid.UUID, status string) (model.PendingUser, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*model.User, error)
}

type AuthService interface {
	Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error)
	Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error)
}

type UUUIDGenerator interface {
	Generate() uuid.UUID
}

type UserService struct {
	log           *slog.Logger
	userStorage   UserStorage
	authService   AuthService
	uuidGenerator UUUIDGenerator
	tracer        trace.Tracer
}

func New(log *slog.Logger,
	userStorage UserStorage,
	authService AuthService,
	uuidGenerator UUUIDGenerator,
	tracer trace.Tracer,
) *UserService {
	return &UserService{
		log:           log,
		userStorage:   userStorage,
		authService:   authService,
		uuidGenerator: uuidGenerator,
		tracer:        tracer,
	}
}

func (u *UserService) processPendingUser(ctx context.Context, userInfo model.UserInfo) (uuid.UUID, error) {
	ctx, span := u.tracer.Start(ctx, "processPendingUser")
	defer span.End()

	pendingUser, err := u.userStorage.PendingUserByEmail(ctx, userInfo.Email)

	if err == nil {
		if pendingUser.Status != model.UserStatusFailed {
			return uuid.Nil, errors.New("wrong pending user status")
		}

		user, err := u.userStorage.ChangePendingUserStatus(ctx, pendingUser.ID, model.UserStatusPending)
		if err != nil {
			return uuid.Nil, err
		}
		return user.ID, nil
	}

	if !errors.Is(err, storage.ErrNotFound) {
		return uuid.Nil, err
	}

	pendingUserInfo := model.PendingUserInfo{
		ID:    u.uuidGenerator.Generate(),
		Email: userInfo.Email,
		Name:  userInfo.Name,
	}
	user, err := u.userStorage.SavePendingUser(ctx, pendingUserInfo)
	if err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}

func (u *UserService) CreateUser(ctx context.Context, userInfo model.UserInfo) error {
	const op = "service.user.CreateUser"
	log := u.log.With(
		slog.String("op", op),
		slog.String("name", userInfo.Email),
	)

	ctx, span := u.tracer.Start(ctx, "CreateUser")
	pendingUserId, err := u.processPendingUser(ctx, userInfo)
	span.End()
	if err != nil {
		log.Error("failed to process pending user", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	ctx, span = u.tracer.Start(ctx, "authService.Register")
	_, err = u.authService.Register(ctx, &ssov1.RegisterRequest{
		Email:    userInfo.Email,
		Password: userInfo.Password,
	})
	span.End()
	if err != nil {
		log.Error("failed to register new user", sl.Err(err))

		_, dbErr := u.userStorage.ChangePendingUserStatus(ctx, pendingUserId, model.UserStatusFailed)
		if dbErr != nil {
			log.Error("failed to change pending user status", sl.Err(err))
		}

		if code, ok := status.FromError(err); ok && code.Code() == codes.AlreadyExists {
			return fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *UserService) Login(ctx context.Context, userInfo model.UserInfo) (string, error) {
	const op = "service.user.Login"
	log := u.log.With(
		slog.String("op", op),
		slog.String("name", userInfo.Email),
	)

	ctx, span := u.tracer.Start(ctx, "authService.Login")
	resp, err := u.authService.Login(ctx, &ssov1.LoginRequest{
		Email:    userInfo.Email,
		Password: userInfo.Password,
		AppId:    "af9956b0-4396-4fc9-8ea5-b00955e87eb6",
	})
	span.End()

	if err != nil {
		log.Error("failed to login user", sl.Err(err))

		code, ok := status.FromError(err)
		if ok {
			switch code.Code() {
			case codes.InvalidArgument:
				return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
			case codes.Unavailable:
				return "", fmt.Errorf("%s: %w", op, ErrUserLocked)
			}
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetToken(), nil
}

func (u *UserService) User(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	const op = "service.user.User"
	log := u.log.With(
		slog.String("op", op),
		slog.String("userId", userID.String()),
	)

	ctx, span := u.tracer.Start(ctx, op)
	defer span.End()

	user, err := u.userStorage.GetUser(ctx, userID)
	if err != nil {
		log.Error("failed retrieving user", sl.Err(err))
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("%s:%w", op, ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return user, nil
}
