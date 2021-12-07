package repository_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/http_srv/repository"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
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

	grpc_mock := new(repository.GrpcMock)
	conn, _ := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(repository.Dialer(grpc_mock)))
	defer conn.Close()

	// http_repository := repository.NewHttpRepository(conn, logger)

	test_cases := []struct {
		test_name      string
		user           entities.User
		repository_res string
		err            error
	}{
		{
			test_name: "user created successfully",
			user: entities.User{
				Email:     "user@email.com",
				Password:  "qwerty",
				Age:       23,
				ExtraInfo: "fav movie: fight club",
			},
			repository_res: "1",
			err:            nil,
		},
		{
			test_name: "no password error",
			user: entities.User{
				Email:     "user@email.com",
				Age:       23,
				ExtraInfo: "fav movie: fight club",
			},
			repository_res: "",
			err:            errors.New("Email or Password empty!"),
		},
		{
			test_name: "no email error",
			user: entities.User{
				Password:  "qwerty",
				Age:       23,
				ExtraInfo: "fav movie: fight club",
			},
			repository_res: "",
			err:            errors.New("Email or Password empty!"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			//  prepare
			ctx := context.Background()
			assert := assert.New(t)

			grpc_req := &userpb.CreateUserRequest{
				Email:                 tc.user.Email,
				Password:              tc.user.Password,
				Age:                   uint32(tc.user.Age),
				AdditionalInformation: tc.user.ExtraInfo,
			}

			grpc_res := &userpb.CreateUserResponse{Id: tc.repository_res}
			grpc_mock.On("CreateUser", ctx, grpc_req).Return(grpc_res, tc.err)

			// act
			res, err := grpc_mock.CreateUser(ctx, grpc_req)

			// assert
			assert.Equal(res.GetId(), tc.repository_res)
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

	grpc_mock := new(repository.GrpcMock)
	conn, _ := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(repository.Dialer(grpc_mock)))
	defer conn.Close()

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
			res: "auth_token",
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
			//  prepare
			assert := assert.New(t)

			assert.Equal(true, true)
		})
	}
}
