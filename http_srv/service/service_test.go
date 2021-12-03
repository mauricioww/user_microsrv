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

type user_data struct {
	Email     string
	Password  string
	Age       int
	ExtraInfo string
}

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
		data      user_data
		res       string
		err       error
	}{
		{
			test_name: "user created successfully",
			data: user_data{
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
			data: user_data{
				Password:  "qwerty",
				Age:       23,
				ExtraInfo: "fav movie: fight club",
			},
			err: errors.New("Email or Password empty!"),
		},
		{
			test_name: "no password error",
			data: user_data{
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
			user := entities.User{
				Email:     tc.data.Email,
				Password:  tc.data.Password,
				Age:       tc.data.Age,
				ExtraInfo: tc.data.ExtraInfo,
			}

			// act
			repository_mock.On("CreateUser", ctx, user).Return(tc.res, tc.err)
			res, err := http_service.CreateUser(ctx, tc.data.Email, tc.data.Password, tc.data.ExtraInfo, tc.data.Age)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(err, tc.err)
		})
	}
}