package userservice

import (
	"context"
	"fmt"
	"testing"

	"github.com/BariVakhidov/rssaggregator/internal/logger"
	"github.com/BariVakhidov/rssaggregator/internal/mocks"
	"github.com/BariVakhidov/rssaggregator/internal/model"
	"github.com/BariVakhidov/rssaggregator/internal/storage"
	ssov1 "github.com/BariVakhidov/ssoprotos/gen/go/sso"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

const passDefaultLen = 10

func TestUserService_CreateUser(t *testing.T) {
	type args struct {
		ctx      context.Context
		userInfo model.UserInfo
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create new user",
			args: args{
				ctx: context.Background(),
				userInfo: model.UserInfo{
					Email:    fmt.Sprintf("test_%s", gofakeit.Email()),
					Name:     gofakeit.LetterN(10),
					Password: generatePassword(),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUUUID := uuid.New()

			pendingUserInfo := getPendingUserInfo(tt.args.userInfo, mockUUUID)

			userStorage := mocks.NewMockUserStorage(t)
			mockUUUIDGenerator := mocks.NewMockUUUIDGenerator(t)
			mockUUUIDGenerator.On("Generate").Once().Return(pendingUserInfo.ID)

			userStorage.On("PendingUserByEmail", mock.Anything, tt.args.userInfo.Email).
				Once().
				Return(model.PendingUser{}, storage.ErrNotFound)
			userStorage.On("SavePendingUser", mock.Anything, pendingUserInfo).
				Once().
				Return(model.PendingUser{
					ID:     pendingUserInfo.ID,
					Email:  pendingUserInfo.Email,
					Name:   pendingUserInfo.Name,
					Status: "pending",
				}, nil)

			authService := mocks.NewMockAuthService(t)
			authService.On("Register", mock.Anything, &ssov1.RegisterRequest{
				Email:    tt.args.userInfo.Email,
				Password: tt.args.userInfo.Password,
			}).
				Once().
				Return(&ssov1.RegisterResponse{}, nil)

			u := &UserService{
				log:           logger.New("local").Log,
				userStorage:   userStorage,
				authService:   authService,
				uuidGenerator: mockUUUIDGenerator,
			}
			if err := u.CreateUser(tt.args.ctx, tt.args.userInfo); (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func generatePassword() string {
	return gofakeit.Password(true, false, true, true, true, passDefaultLen)
}

func getPendingUserInfo(userInfo model.UserInfo, mockUUUID uuid.UUID) model.PendingUserInfo {
	return model.PendingUserInfo{
		ID:    mockUUUID,
		Email: userInfo.Email,
		Name:  userInfo.Name,
	}
}
