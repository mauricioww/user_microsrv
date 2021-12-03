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
