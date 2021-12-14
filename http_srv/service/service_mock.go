package service

import (
	"context"

	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/stretchr/testify/mock"
)

type RepoMock struct {
	mock.Mock
}

func (r *RepoMock) CreateUser(ctx context.Context, user entities.User) (string, error) {
	args := r.Called(ctx, user)

	return args.String(0), args.Error(1)
}

func (r *RepoMock) Authenticate(ctx context.Context, session entities.Session) (int, error) {
	args := r.Called(ctx, session)

	return args.Int(0), args.Error(1)
}

func (r *RepoMock) UpdateUser(ctx context.Context, user entities.UserUpdate) (entities.User, error) {
	args := r.Called(ctx, user)

	return args.Get(0).(entities.User), args.Error(1)
}
