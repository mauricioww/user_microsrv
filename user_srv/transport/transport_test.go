package transport_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/mauricioww/user_microsrv/user_srv/transport"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	mock_srv := new(transport.GrpcUserSrvMock)
	endpoints := transport.MakeGrpcUserServiceEndpoints(mock_srv)
	grpc_service := transport.NewGrpcUserServer(endpoints)

	test_cases := []struct {
		test_name string
		user_req  *userpb.CreateUserRequest
		user_res  *userpb.CreateUserResponse
		srv_res   int
		err       error
	}{
		{
			test_name: "user created successfully",
			user_req: &userpb.CreateUserRequest{
				Email:    "success@email.com",
				Password: "qwerty",
				Age:      23,
			},
			srv_res: 0,
			err:     nil,
		},
		{
			test_name: "no password error",
			user_req: &userpb.CreateUserRequest{
				Email: "success@email.com",
				Age:   23,
			},
			srv_res: -1,
			err:     errors.New("Password empty"),
		},
		{
			test_name: "no email error",
			user_req: &userpb.CreateUserRequest{
				Password: "qwerty",
				Age:      23,
			},
			srv_res: -1,
			err:     errors.New("Email empty"),
		},
	}
	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()
			if tc.err != nil {
				tc.user_res = nil
			} else {
				tc.user_res = &userpb.CreateUserResponse{Id: int32(tc.srv_res)}
			}

			// act
			mock_srv.On("CreateUser", ctx, tc.user_req.GetEmail(), tc.user_req.GetPassword(), int(tc.user_req.GetAge())).Return(tc.srv_res, tc.err)
			res, err := grpc_service.CreateUser(ctx, tc.user_req)

			// assert
			assert.Equal(tc.user_res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestAuthenticate(t *testing.T) {
	mock_srv := new(transport.GrpcUserSrvMock)
	endpoints := transport.MakeGrpcUserServiceEndpoints(mock_srv)
	grpc_service := transport.NewGrpcUserServer(endpoints)

	test_cases := []struct {
		test_name string
		data      *userpb.AuthenticateRequest
		res       *userpb.AuthenticateResponse
		srv_res   int
		err       error
	}{
		{
			test_name: "authenticate successfully",
			data: &userpb.AuthenticateRequest{
				Email:    "user@email.com",
				Password: "qwerty",
			},
			err: nil,
		},
		{
			test_name: "no password error",
			data: &userpb.AuthenticateRequest{
				Email: "user@email.com",
			},
			srv_res: -1,
			err:     errors.New("Password empty"),
		},
		{
			test_name: "no email error",
			data: &userpb.AuthenticateRequest{
				Password: "invalid_password",
			},
			srv_res: -1,
			err:     errors.New("Email empty"),
		},
		{
			test_name: "invalid password error",
			data: &userpb.AuthenticateRequest{
				Email:    "user@email.com",
				Password: "invalid_password",
			},
			srv_res: -1,
			err:     errors.New("Password error"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()
			if tc.err != nil {
				tc.res = nil
			} else {
				tc.res = &userpb.AuthenticateResponse{UserId: int32(tc.srv_res)}
			}

			// act
			mock_srv.On("Authenticate", ctx, tc.data.GetEmail(), tc.data.GetPassword()).Return(tc.srv_res, tc.err)
			res, err := grpc_service.Authenticate(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	mock_srv := new(transport.GrpcUserSrvMock)
	endpoints := transport.MakeGrpcUserServiceEndpoints(mock_srv)
	grpc_service := transport.NewGrpcUserServer(endpoints)

	test_cases := []struct {
		test_name string
		data      *userpb.UpdateUserRequest
		res       *userpb.UpdateUserResponse
		srv_res   bool
		err       error
	}{
		{
			test_name: "update user successfully",
			data: &userpb.UpdateUserRequest{
				Id:       1,
				Email:    "new_email@domain.com",
				Password: "new_password",
				Age:      25,
			},
			srv_res: true,
			err:     nil,
		},
		{
			test_name: "no password error",
			data: &userpb.UpdateUserRequest{
				Id:    1,
				Email: "new_email@domain.com",
				Age:   25,
			},
			err: errors.New("Password empty"),
		},
		{
			test_name: "no email error",
			data: &userpb.UpdateUserRequest{
				Id:       1,
				Password: "new_password",
				Age:      25,
			},
			err: errors.New("Email empty"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()
			if tc.err != nil {
				tc.res = nil
			} else {
				tc.res = &userpb.UpdateUserResponse{Success: tc.srv_res}
			}

			// act
			mock_srv.On("UpdateUser", ctx, int(tc.data.GetId()), tc.data.GetEmail(), tc.data.GetPassword(), int(tc.data.GetAge())).Return(tc.srv_res, tc.err)
			res, err := grpc_service.UpdateUser(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestGetUser(t *testing.T) {
	mock_srv := new(transport.GrpcUserSrvMock)
	endpoints := transport.MakeGrpcUserServiceEndpoints(mock_srv)
	grpc_service := transport.NewGrpcUserServer(endpoints)

	test_cases := []struct {
		test_name string
		data      *userpb.GetUserRequest
		res       *userpb.GetUserResponse
		srv_res   entities.User
		err       error
	}{
		{
			test_name: "user found",
			data: &userpb.GetUserRequest{
				Id: 0,
			},
			srv_res: entities.User{
				Email:    "user@email.com",
				Password: "password",
				Age:      20,
			},
			err: nil,
		},
		{
			test_name: "user not found error",
			data: &userpb.GetUserRequest{
				Id: 1,
			},
			err: errors.New("User not found"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()
			if tc.err != nil {
				tc.res = nil
			} else {
				tc.res = &userpb.GetUserResponse{
					Email:    tc.srv_res.Email,
					Password: tc.srv_res.Password,
					Age:      uint32(tc.srv_res.Age),
				}
			}

			// act
			mock_srv.On("GetUser", ctx, int(tc.data.GetId())).Return(tc.srv_res, tc.err)
			res, err := grpc_service.GetUser(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	mock_srv := new(transport.GrpcUserSrvMock)
	endpoints := transport.MakeGrpcUserServiceEndpoints(mock_srv)
	grpc_service := transport.NewGrpcUserServer(endpoints)

	test_cases := []struct {
		test_name string
		data      *userpb.DeleteUserRequest
		res       *userpb.DeleteUserResponse
		srv_res   bool
		err       error
	}{
		{
			test_name: "delete user success",
			data: &userpb.DeleteUserRequest{
				Id: 0,
			},
			srv_res: true,
			err:     nil,
		},
		{
			test_name: "user not found error",
			data: &userpb.DeleteUserRequest{
				Id: 1,
			},
			err: errors.New("User not found"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()
			if tc.err != nil {
				tc.res = nil
			} else {
				tc.res = &userpb.DeleteUserResponse{Success: tc.srv_res}
			}

			// act
			mock_srv.On("DeleteUser", ctx, int(tc.data.GetId())).Return(tc.srv_res, tc.err)
			res, err := grpc_service.DeleteUser(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}
