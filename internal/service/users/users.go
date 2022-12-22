package users

import (
	"context"
	"github.com/google/uuid"
	usersDomain "homework/internal/domain/users"
)

type service struct {
	usersStorage UsersStorage
}

type UsersService interface {
	GetUserById(ctx context.Context, userId uuid.UUID) (*usersDomain.User, error)
}

type UsersStorage interface {
	GetUserById(ctx context.Context, userId uuid.UUID) (*usersDomain.User, error)
}

func (s service) GetUserById(ctx context.Context, userId uuid.UUID) (*usersDomain.User, error) {
	return s.usersStorage.GetUserById(ctx, userId)
}

func NewUsersService(usersStorage UsersStorage) UsersService {
	return &service{usersStorage: usersStorage}
}
