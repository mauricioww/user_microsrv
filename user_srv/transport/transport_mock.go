package transport

import (
	"context"

	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/stretchr/testify/mock"
)

type GrpcUserSrvMock struct {
	mock.Mock
}

func (s *GrpcUserSrvMock) CreateUser(ctx context.Context, email string, pwd string, extra_info string, age int) (string, error) {
	args := s.Called(ctx, email, pwd, extra_info, age)

	return args.String(0), args.Error(1)
}

func (s *GrpcUserSrvMock) Authenticate(ctx context.Context, email string, pwd string) (int, error) {
	args := s.Called(ctx, email, pwd)

	return args.Int(0), args.Error(1)
}

func (s *GrpcUserSrvMock) UpdateUser(ctx context.Context, id int, email string, pwd string, extra_info string, age int) (entities.User, error) {
	args := s.Called(ctx, id, email, pwd, extra_info, age)

	return args.Get(0).(entities.User), args.Error(1)
}

func (s *GrpcUserSrvMock) GetUser(ctx context.Context, id int) (entities.User, error) {
	args := s.Called(ctx, id)

	return args.Get(0).(entities.User), args.Error(1)
}
