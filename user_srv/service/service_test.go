package service_test

import (
	"context"
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
	grpc_user_srv = service.NewGrpcUserService(logger, user_repo_mock)

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
