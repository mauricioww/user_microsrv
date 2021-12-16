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

func (r *UserRepositoryMock) Authenticate(ctx context.Context, session *entities.Session) (string, error) {
	args := r.Called(ctx, session)

	return args.String(0), args.Error(1)
}

func (r *UserRepositoryMock) UpdateUser(ctx context.Context, update entities.Update) (entities.User, error) {
	args := r.Called(ctx, update)

	return args.Get(0).(entities.User), args.Error(1)
}

func (r *UserRepositoryMock) GetUser(ctx context.Context, id int) (entities.User, error) {
	args := r.Called(ctx, id)

	return args.Get(0).(entities.User), args.Error(1)
}
