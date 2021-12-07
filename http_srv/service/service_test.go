package service_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/http_srv/service"
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

	repository_mock := new(service.RepoMock)
	http_service := service.NewHttpService(repository_mock, logger)

	test_cases := []struct {
		test_name string
		data      entities.User
		res       string
		err       error
	}{
		{
			test_name: "user created successfully",
			data: entities.User{
				Email:     "success@email.com",
				Password:  "qwerty",
				Age:       23,
				ExtraInfo: "fav movie: fight club",
			},
			res: "success@email.com",
			err: nil,
		},
		{
			test_name: "no email error",
			data: entities.User{
				Password:  "qwerty",
				Age:       23,
				ExtraInfo: "fav movie: fight club",
			},
			err: errors.New("Email or Password empty!"),
		},
		{
			test_name: "no password error",
			data: entities.User{
				Email:     "success@email.com",
				Age:       23,
				ExtraInfo: "fav movie: fight club",
			},
			err: errors.New("Email or Password empty!"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			repository_mock.On("CreateUser", ctx, tc.data).Return(tc.res, tc.err)
			res, err := http_service.CreateUser(ctx, tc.data.Email, tc.data.Password, tc.data.ExtraInfo, tc.data.Age)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(err, tc.err)
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

	repository_mock := new(service.RepoMock)
	http_service := service.NewHttpService(repository_mock, logger)

	test_cases := []struct {
		test_name string
		data      entities.Session
		res       string
		err       error
	}{
		{
			test_name: "success authenticate",
			data: entities.Session{
				Email:    "fake_email@email.com",
				Password: "fake_password",
			},
			res: "generated_auth_token",
			err: nil,
		},
		{
			test_name: "no email error",
			data: entities.Session{
				Password: "fake_password",
			},
			err: errors.New("Email or Password empty!"),
		},
		{
			test_name: "no password error",
			data: entities.Session{
				Email: "fake_email@email.com",
			},
			err: errors.New("Email or Password empty!"),
		},
		{
			test_name: "user not found",
			data: entities.Session{
				Email:    "no_real@email.com",
				Password: "fake_password",
			},
			res: "",
			err: errors.New("User not found"),
		},
		{
			test_name: "invalid password",
			data: entities.Session{
				Email:    "user@email.com",
				Password: "invalid_password",
			},
			res: "",
			err: errors.New("Invalid user"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			repository_mock.On("Authenticate", ctx, tc.data).Return(tc.res, tc.err)
			res, err := http_service.Authenticate(ctx, tc.data.Email, tc.data.Password)

			// assert
			if tc.res == "" {
				assert.Empty(res)
			} else {
				assert.NotEmpty(res)
			}
			assert.Equal(tc.err, err)
		})
	}
}
