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
		repository_res int
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
			repository_res: 1,
			err:            nil,
		},
		{
			test_name: "no password error",
			user: entities.User{
				Email:     "user@email.com",
				Age:       23,
				ExtraInfo: "fav movie: fight club",
			},
			repository_res: -1,
			err:            errors.New("Email or Password empty!"),
		},
		{
			test_name: "no email error",
			user: entities.User{
				Password:  "qwerty",
				Age:       23,
				ExtraInfo: "fav movie: fight club",
			},
			repository_res: -1,
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

			grpc_res := &userpb.CreateUserResponse{Id: int32(tc.repository_res)}
			grpc_mock.On("CreateUser", ctx, grpc_req).Return(grpc_res, tc.err)

			// act
			res, err := grpc_mock.CreateUser(ctx, grpc_req)

			// assert
			assert.Equal(int(res.GetId()), tc.repository_res)
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
		res       int
		err       error
	}{
		{
			test_name: "success authenticate",
			data: entities.Session{
				Email:    "fake_email@email.com",
				Password: "fake_password",
			},
			res: 1,
			err: nil,
		},
		{
			test_name: "no email error",
			data: entities.Session{
				Password: "fake_password",
			},
			res: -1,
			err: errors.New("Email or Password empty!"),
		},
		{
			test_name: "no password error",
			data: entities.Session{
				Email: "fake_email@email.com",
			},
			res: -1,
			err: errors.New("Email or Password empty!"),
		},
		{
			test_name: "user not found",
			data: entities.Session{
				Email:    "no_real@email.com",
				Password: "fake_password",
			},
			res: -1,
			err: errors.New("User not found"),
		},
		{
			test_name: "invalid password",
			data: entities.Session{
				Email:    "user@email.com",
				Password: "invalid_password",
			},
			res: -1,
			err: errors.New("Invalid user"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			//  prepare
			assert := assert.New(t)

			ctx := context.Background()

			grpc_req := &userpb.AuthenticateRequest{
				Email:    tc.data.Email,
				Password: tc.data.Password,
			}

			grpc_res := &userpb.AuthenticateResponse{UserId: int32(tc.res)}
			grpc_mock.On("Authenticate", ctx, grpc_req).Return(grpc_res, tc.err)

			// act
			res, err := grpc_mock.Authenticate(ctx, grpc_req)

			// assert
			assert.Equal(res.GetUserId(), int32(tc.res))
			assert.Equal(err, tc.err)
		})
	}
}

func UpdateUser(t *testing.T) {
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
		test_name string
		fields    []string
		prev_user entities.User
		aft_user  entities.User
		err       error
	}{
		{
			test_name: "Update email successfully",
			prev_user: entities.User{
				Email: "user@email.com",
			},
			aft_user: entities.User{
				Email: "new_email@domain.com",
			},
			err: nil,
		},
		{
			test_name: "Update email and password successfully",
			prev_user: entities.User{
				Email:    "user@email.com",
				Password: "password",
			},
			aft_user: entities.User{
				Email:    "new_email@domain.com",
				Password: "new_password",
			},
			err: nil,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			//  prepare
			ctx := context.Background()
			assert := assert.New(t)

			grpc_req := &userpb.UpdateUserRequest{
				Email:                 tc.prev_user.Email,
				Password:              tc.prev_user.Password,
				Age:                   uint32(tc.prev_user.Age),
				AdditionalInformation: tc.prev_user.ExtraInfo,
			}

			grpc_res := &userpb.UpdateUserResponse{
				Email:                 tc.aft_user.Email,
				Password:              tc.aft_user.Password,
				Age:                   uint32(tc.aft_user.Age),
				AdditionalInformation: tc.aft_user.ExtraInfo,
			}
			grpc_mock.On("UpdateUser", ctx, grpc_req).Return(grpc_res, tc.err)

			// act
			res, err := grpc_mock.UpdateUser(ctx, grpc_req)

			// assert
			assert.Equal(res, grpc_res)
			assert.Equal(err, tc.err)
		})
	}
}

func TestGetUser(t *testing.T) {
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
		data      int
		res       entities.User
		err       error
	}{
		{
			test_name: "user found",
			data:      1,
			res: entities.User{
				Email:     "email@domain.com",
				Password:  "password",
				Age:       10,
				ExtraInfo: "extra_info",
			},
			err: nil,
		},
		{
			test_name: "user not found",
			data:      -1,
			res:       entities.User{},
			err:       errors.New("User not found"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			//  prepare
			assert := assert.New(t)

			ctx := context.Background()

			grpc_req := &userpb.GetUserRequest{
				Id: uint32(tc.data),
			}

			grpc_res := &userpb.GetUserResponse{
				Email:                 tc.res.Email,
				Password:              tc.res.Password,
				Age:                   uint32(tc.res.Age),
				AdditionalInformation: tc.res.ExtraInfo,
			}
			grpc_mock.On("GetUser", ctx, grpc_req).Return(grpc_res, tc.err)

			// act
			res, err := grpc_mock.GetUser(ctx, grpc_req)

			// assert
			assert.Equal(res, grpc_res)
			assert.Equal(err, tc.err)
		})
	}
}

func TestDeleteUser(t *testing.T) {
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
		test_name string
		data      int
		res       bool
		err       error
	}{
		{
			test_name: "user deleted success",
			data:      1,
			res:       true,
			err:       nil,
		},
		{
			test_name: "user delete error",
			data:      -1,
			res:       false,
			err:       errors.New("User not found"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			//  prepare
			assert := assert.New(t)

			ctx := context.Background()

			grpc_req := &userpb.DeleteUserRequest{Id: uint32(tc.data)}

			grpc_res := &userpb.DeleteUserResponse{Success: tc.res}

			grpc_mock.On("DeleteUser", ctx, grpc_req).Return(grpc_res, tc.err)

			// act
			res, err := grpc_mock.DeleteUser(ctx, grpc_req)

			// assert
			assert.Equal(res, grpc_res)
			assert.Equal(err, tc.err)
		})
	}
}
