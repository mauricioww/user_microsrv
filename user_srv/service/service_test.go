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
	"github.com/stretchr/testify/mock"
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
			user_repo_mock.On("CreateUser", ctx, mock.AnythingOfType("entities.User")).Return(tc.res, tc.err)
			_, err := grpc_user_srv.CreateUser(ctx, tc.data.Email, tc.data.Password, tc.data.ExtraInfo, tc.data.Age)

			// assert
			// assert.Equal(tc.res, res)
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
		data      *entities.Session
		repo_pwd  string
		repo_err  error
		res       int
		err       error
	}{
		{
			test_name: "authenticate successfully",
			data: &entities.Session{
				Email:    "user@email.com",
				Password: "secret",
			},
			repo_pwd: "secret",
		},
		{
			test_name: "user not found error",
			data: &entities.Session{
				Email:    "fake_user@email.com",
				Password: "fake_password",
			},
			repo_pwd: "fake_password",
			repo_err: errors.New("User not found"),
			res:      -1,
			err:      errors.New("User not found"),
		},
		{
			test_name: "invalid pasword error",
			data: &entities.Session{
				Email:    "user@email.com",
				Password: "fakee_password",
			},
			repo_pwd: "incorrect_password",
			res:      -1,
			err:      errors.New("Password error"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)
			// invalid data

			// act
			user_repo_mock.On("Authenticate", ctx, tc.data).Return(tc.repo_pwd, tc.repo_err)
			res, err := grpc_user_srv.Authenticate(ctx, tc.data.Email, tc.data.Password)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
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
		data      entities.Update
		res       entities.User
		err       error
	}{
		{
			test_name: "update any field",
			data: entities.Update{
				UserId: 0,
				Information: map[string]interface{}{
					"email": "new_email@domain.com",
				},
			},
			res: entities.User{
				Email: "new_email@domain.com",
			},
			err: nil,
		},
		{
			test_name: "update one string and integer field",
			data: entities.Update{
				UserId: 0,
				Information: map[string]interface{}{
					"password": "new_password",
					"age":      23,
				},
			},
			res: entities.User{
				Password: "new_password",
				Age:      23,
			},
			err: nil,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			user_repo_mock.On("UpdateUser", ctx, tc.data).Return(tc.res, tc.err)
			res, err := grpc_user_srv.UpdateUser(ctx, tc.data.UserId, tc.data.Information)

			// assert
			assert.Equal(tc.err, err)
			assert.Equal(tc.res, res)
		})
	}
}
