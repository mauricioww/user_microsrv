package transport_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/http_srv/transport"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	http_srv_mock := new(transport.ServiceMock)

	endpoints := transport.MakeHttpEndpoints(http_srv_mock)

	test_cases := []struct {
		test_name  string
		res_string int
		data       transport.CreateUserRequest
		res        transport.CreateUserResponse
		err        error
	}{
		{
			test_name: "user created successfully",
			data: transport.CreateUserRequest{
				Email:    "success@email.com",
				Password: "qwerty",
				Age:      23,
				Details:  transport.GenereateDetails(),
			},
			res_string: 1,
			res: transport.CreateUserResponse{
				Id:       1,
				Email:    "success@email.com",
				Password: "qwerty",
				Age:      23,
				Details:  transport.GenereateDetails(),
			},
			err: nil,
		},
		{
			test_name: "no password error",
			data: transport.CreateUserRequest{
				Email: "success@email.com",
				Age:   23,
			},
			res_string: -1,
			res: transport.CreateUserResponse{
				Id:    -1,
				Email: "success@email.com",
				Age:   23,
			},
			err: errors.New("Email or Password empty!"),
		},
		{
			test_name: "no email error",
			data: transport.CreateUserRequest{
				Password: "qwerty",
				Age:      23,
			},
			res_string: -1,
			res: transport.CreateUserResponse{
				Id:       -1,
				Password: "qwerty",
				Age:      23,
			},
			err: errors.New("Email or Password empty!"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()

			http_srv_mock.On("CreateUser", ctx, tc.data.Email, tc.data.Password, tc.data.Age, tc.data.Details).
				Return(tc.res_string, tc.err)

			// act
			res, err := endpoints.CreateUser(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestAuthenticate(t *testing.T) {

	http_srv_mock := new(transport.ServiceMock)

	endpoints := transport.MakeHttpEndpoints(http_srv_mock)

	test_cases := []struct {
		test_name  string
		data       transport.AuthenticateRequest
		res        transport.AuthenticateResponse
		res_string string
		err        error
	}{
		{
			test_name: "success authenticate",
			data: transport.AuthenticateRequest{
				Email:    "user@email.com",
				Password: "fake_password",
			},
			res: transport.AuthenticateResponse{
				Token: "auth_token",
			},
			res_string: "auth_token",
			err:        nil,
		},
		{
			test_name: "no password error",
			data: transport.AuthenticateRequest{
				Email: "user@email.com",
			},
			err: errors.New("Empty email or password!"),
		},
		{
			test_name: "no email error",
			data: transport.AuthenticateRequest{
				Password: "fake_password",
			},
			err: errors.New("Empty email or password!"),
		},
		{
			test_name: "user not found",
			data: transport.AuthenticateRequest{
				Email:    "no_real@email.com",
				Password: "fake_password",
			},
			res_string: "",
			err:        errors.New("User not found"),
		},
		{
			test_name: "invalid password",
			data: transport.AuthenticateRequest{
				Email:    "user@email.com",
				Password: "invalid_password",
			},
			res_string: "",
			err:        errors.New("Invalid user"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			http_srv_mock.On("Authenticate", ctx, tc.data.Email, tc.data.Password).Return(tc.res_string, tc.err)
			// act
			res, err := endpoints.Authenticate(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	http_srv_mock := new(transport.ServiceMock)

	endpoints := transport.MakeHttpEndpoints(http_srv_mock)

	test_cases := []struct {
		test_name string
		data      transport.UpdateUserRequest
		res       transport.UpdateUserResponse
		user_res  entities.User
		err       error
	}{
		{
			test_name: "update user success",
			data: transport.UpdateUserRequest{
				UserId:   1,
				Email:    "new_email@domain.com",
				Password: "new_password",
				Age:      23,
			},
			user_res: entities.User{
				Email:    "new_email@domain.com",
				Password: "new_password",
				Age:      23,
			},
			res: transport.UpdateUserResponse{
				Id:       1,
				Email:    "new_email@domain.com",
				Password: "new_password",
				Age:      23,
			},
		},
	}

	for _, tc := range test_cases {
		// prepare
		ctx := context.Background()
		assert := assert.New(t)

		http_srv_mock.On("UpdateUser", ctx, tc.data.UserId, tc.data.Email, tc.data.Password, tc.data.Age).Return(tc.user_res, tc.err)

		// act
		res, err := endpoints.UpdateUser(ctx, tc.data)

		// assert
		assert.Equal(tc.res, res)
		assert.Equal(tc.err, err)
	}
}

func TestGetUser(t *testing.T) {
	http_srv_mock := new(transport.ServiceMock)

	endpoints := transport.MakeHttpEndpoints(http_srv_mock)

	test_cases := []struct {
		test_name string
		data      transport.GetUserRequest
		res       transport.GetUserResponse
		user_res  entities.User
		err       error
	}{
		{
			test_name: "user found success",
			data: transport.GetUserRequest{
				UserId: 1,
			},
			user_res: entities.User{
				Email:    "email@email.com",
				Password: "password",
				Age:      23,
			},
			res: transport.GetUserResponse{
				Id:       1,
				Email:    "email@email.com",
				Password: "password",
				Age:      23,
			},
		},
	}

	for _, tc := range test_cases {
		// prepare
		ctx := context.Background()
		assert := assert.New(t)

		http_srv_mock.On("GetUser", ctx, tc.data.UserId).Return(tc.user_res, tc.err)

		// act
		res, err := endpoints.GetUser(ctx, tc.data)

		// assert
		assert.Equal(tc.res, res)
		assert.Equal(tc.err, err)
	}
}
