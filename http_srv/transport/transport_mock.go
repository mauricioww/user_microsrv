package transport

import (
	"context"

	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) CreateUser(ctx context.Context, email string, pwd string, extra_info string, age int) (int, error) {
	args := s.Called(ctx, email, pwd, extra_info, age)

	return args.Int(0), args.Error(1)
}

func (s *ServiceMock) Authenticate(ctx context.Context, email string, pwd string) (string, error) {
	args := s.Called(ctx, email, pwd)

	return args.String(0), args.Error(1)
}

func (s *ServiceMock) UpdateUser(ctx context.Context, user_id int, email string, pwd string, extra_info string, age int) (entities.User, error) {
	args := s.Called(ctx, user_id, email, pwd, extra_info, age)

	return args.Get(0).(entities.User), args.Error(1)
}

func (s *ServiceMock) GetUser(ctx context.Context, user_id int) (entities.User, error) {
	args := s.Called(ctx, user_id)

	return args.Get(0).(entities.User), args.Error(1)
}

func (s *ServiceMock) DeleteUser(ctx context.Context, user_id int) (bool, error) {
	args := s.Called(ctx, user_id)

	return args.Bool(0), args.Error(1)
}
