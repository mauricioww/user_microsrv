package service

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type RepoMock struct {
	mock.Mock
}

func (r *RepoMock) CreateUser(ctx context.Context, email string, pwd string, extra_info string, age int) (string, error) {
	args := r.Called(ctx, email, pwd, extra_info, age)

	return args.Get(0).(string), args.Error(1)
}
