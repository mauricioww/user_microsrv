package service

import (
	"context"

	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (r *UserRepositoryMock) CreateUser(ctx context.Context, user entities.User) (string, error) {
	args := r.Called(ctx, user)

	return args.String(0), args.Error(1)
}
