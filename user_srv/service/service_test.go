package service_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/mauricioww/user_microsrv/user_srv/service"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(
			logger,
			"service",
			"account",
			"time",
			log.DefaultTimestampUTC,
			"caller",
			log.DefaultCaller,
		)
	}

	var grpc_user_srv service.GrpcUserService

	user_repo_mock := new(service.UserRepositoryMock)
	grpc_user_srv = service.NewGrpcUserService(user_repo_mock, logger)

	test_cases := []struct {
		test_name string
		data      entities.User
		res       string
		err       error
	}{
		{
			test_name: "create user successfully",
			data: entities.User{
				Email:     "success@email.com",
				Password:  "qwerty",
				Age:       23,
				ExtraInfo: "fav movie: fight club",
			},
			res: "1",
			err: nil,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			user_repo_mock.On("CreateUser", ctx, tc.data).Return(tc.res, tc.err)
			res, err := grpc_user_srv.CreateUser(ctx, tc.data.Email, tc.data.Password, tc.data.ExtraInfo, tc.data.Age)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestAuthenticate(t *testing.T) {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(
			logger,
			"service",
			"account",
			"time",
			log.DefaultTimestampUTC,
			"caller",
			log.DefaultCaller,
		)
	}

	var grpc_user_srv service.GrpcUserService

	user_repo_mock := new(service.UserRepositoryMock)
	grpc_user_srv = service.NewGrpcUserService(user_repo_mock, logger)

	test_cases := []struct {
		test_name string
		data      entities.Session
		res       string
		err       error
	}{
		{
			test_name: "authenticate successfully",
			data: entities.Session{
				Email:    "user@email.com",
				Password: "fake_password",
			},
			res: "user_authenticated",
			err: nil,
		},
		{
			test_name: "user not found",
			data: entities.Session{
				Email:    "fake_user@email.com",
				Password: "fake_password",
			},
			res: "",
			err: errors.New("Invalid password"),
		},
		{
			test_name: "invalid pasword",
			data: entities.Session{
				Email:    "user@email.com",
				Password: "invalid_password",
			},
			res: "",
			err: errors.New("User not found"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			user_repo_mock.On("Authenticate", ctx, tc.data).Return(tc.data.Password, tc.err)
			res, err := grpc_user_srv.Authenticate(ctx, tc.data.Email, tc.data.Password)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}
