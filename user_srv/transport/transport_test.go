package transport_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/mauricioww/user_microsrv/user_srv/transport"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	grpc_user_srv_mock := new(transport.GrpcUserSrvMock)

	endpoints := transport.MakeGrpcUserServiceEndpoints(grpc_user_srv_mock)

	test_cases := []struct {
		test_name string
		data      transport.CreateUserRequest
		res       string
		err       error
	}{
		{
			test_name: "user created successfully",
			data: transport.CreateUserRequest{
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
			assert := assert.New(t)
			ctx := context.Background()

			grpc_user_srv_mock.On("CreateUser", ctx, tc.data.Email, tc.data.Password, tc.data.ExtraInfo, tc.data.Age).Return(tc.res, tc.err)

			// act
			res, err := endpoints.CreateUser(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestAuthenticate(t *testing.T) {
	grpc_user_srv_mock := new(transport.GrpcUserSrvMock)

	endpoints := transport.MakeGrpcUserServiceEndpoints(grpc_user_srv_mock)

	test_cases := []struct {
		test_name string
		data      transport.AuthenticateRequest
		res       int
		err       error
	}{
		{
			test_name: "authenticate successfully",
			data: transport.AuthenticateRequest{
				Email:    "user@email.com",
				Password: "qwerty",
			},
			res: 0,
			err: nil,
		},
		{
			test_name: "invalid password error",
			data: transport.AuthenticateRequest{
				Email:    "user@email.com",
				Password: "invalid_password",
			},
			res: -1,
			err: errors.New("Invalid Password"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()

			grpc_user_srv_mock.On("Authenticate", ctx, tc.data.Email, tc.data.Password).Return(tc.res, tc.err)

			// act
			res, err := endpoints.Authenticate(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	grpc_user_srv_mock := new(transport.GrpcUserSrvMock)

	endpoints := transport.MakeGrpcUserServiceEndpoints(grpc_user_srv_mock)

	test_cases := []struct {
		test_name string
		data      transport.UpdateUserRequest
		res       entities.User
		err       error
	}{
		{
			test_name: "update user successfully",
			data: transport.UpdateUserRequest{
				Id:       1,
				Email:    "new_email@domain.com",
				Password: "new_password",
				Age:      25,
			},
			res: entities.User{
				Email:     "new_email@domain.com",
				Password:  "new_password",
				Age:       25,
				ExtraInfo: "",
			},
			err: nil,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()

			grpc_user_srv_mock.On("UpdateUser", ctx, tc.data.Id, tc.data.Email, tc.data.Password, tc.data.ExtraInfo, tc.data.Age).Return(tc.res, tc.err)

			// act
			res, err := endpoints.UpdateUser(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}
