package service

import (
	"context"

	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/stretchr/testify/mock"
)

type RepoMock struct {
	mock.Mock
}

func (r *RepoMock) CreateUser(ctx context.Context, user entities.User) (int, error) {
	args := r.Called(ctx, user)

	return args.Int(0), args.Error(1)
}

func (r *RepoMock) Authenticate(ctx context.Context, session entities.Session) (int, error) {
	args := r.Called(ctx, session)

	return args.Int(0), args.Error(1)
}

func (r *RepoMock) UpdateUser(ctx context.Context, user entities.UserUpdate) (bool, error) {
	args := r.Called(ctx, user)

	return args.Bool(0), args.Error(1)
}

func (r *RepoMock) GetUser(ctx context.Context, id int) (entities.User, error) {
	args := r.Called(ctx, id)

	return args.Get(0).(entities.User), args.Error(1)
}

func (r *RepoMock) DeleteUser(ctx context.Context, id int) (bool, error) {
	args := r.Called(ctx, id)

	return args.Bool(0), args.Error(1)
}

func GenenerateDetails() entities.Details {
	return entities.Details{
		Country:      "Mexico",
		City:         "CDMX",
		MobileNumber: "11223344",
		Married:      false,
		Height:       1.75,
		Weigth:       76.0,
	}
}
