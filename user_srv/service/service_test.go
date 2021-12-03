package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/mauricioww/user_microsrv/user_srv/service"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
	"github.com/stretchr/testify/assert"
)

func CreateUserTest(t *testing.T) {
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
		data      *userpb.CreateUserRequest
		res       *userpb.CreateUserResponse
		err       error
	}{
		{
			test_name: "create user successfully",
			data: &userpb.CreateUserRequest{
				Email:                 "success@email.com",
				Password:              "qwerty",
				Age:                   23,
				AdditionalInformation: "fav movie: fight club",
			},
			res: &userpb.CreateUserResponse{
				Id: "1",
			},
			err: nil,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)
			user := entities.User{
				Email:     tc.data.GetEmail(),
				Password:  tc.data.GetPassword(),
				Age:       int(tc.data.GetAge()),
				ExtraInfo: tc.data.GetAdditionalInformation(),
			}

			// act
			user_repo_mock.On("CreateUser", ctx, user).Return(tc.res, tc.err)
			res, err := grpc_user_srv.CreateUser(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}
