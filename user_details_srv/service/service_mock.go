package service

import (
	"context"

	"github.com/mauricioww/user_microsrv/user_details_srv/entities"
	"github.com/stretchr/testify/mock"
)

type UserDetailsRepositoryMock struct {
	mock.Mock
}

func (r *UserDetailsRepositoryMock) SetUserDetails(ctx context.Context, information entities.UserDetails) (bool, error) {
	args := r.Called(ctx, information)

	return args.Bool(0), args.Error(1)
}

func (r *UserDetailsRepositoryMock) GetUserDetails(ctx context.Context, user_id int) (entities.UserDetails, error) {
	args := r.Called(ctx, user_id)

	return args.Get(0).(entities.UserDetails), args.Error(1)
}
