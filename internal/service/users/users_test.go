package users

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	usersDomain "homework/internal/domain/users"
	mockUsersService "homework/internal/service/users/mock"
	"homework/internal/util/terr"
)

//go:generate mockgen -destination ./mock/users_service_mock.go homework/internal/service/users UsersService

func Test_GetUserByID(t *testing.T) {

	// Arrange
	userId := uuid.MustParse("244f9f9a-f730-4860-b5aa-479c19320fa5")

	var tests = []struct {
		name string
		args uuid.UUID
		want *usersDomain.User
		err  error
	}{
		{
			name: "success",
			args: userId,
			want: &usersDomain.User{
				Id:    userId,
				Name:  "User 123",
				Email: "123@gmail.com",
			},
			err: nil,
		},
		{
			name: "fail/user not found",
			args: userId,
			want: nil,
			err:  terr.NotFound(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			usersService := mockUsersService.NewMockUsersService(ctrl)
			usersService.EXPECT().
				GetUserByID(ctx, tt.args).
				Return(tt.want, tt.err)

			// Act
			got, err := usersService.GetUserByID(ctx, tt.args)

			// Assert
			assert.Equal(t, tt.err, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}

}
