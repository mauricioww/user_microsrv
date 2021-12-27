package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/http_srv/repository"
	"github.com/mauricioww/user_microsrv/user_details_srv/detailspb"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {

	user_mock := new(repository.GrpcUserMock)
	details_mock := new(repository.GrpcDetailsMock)
	conn1, conn2, http_repository := repository.InitRepoMock(user_mock, details_mock)

	defer conn1.Close()
	defer conn2.Close()

	test_cases := []struct {
		test_name      string
		user           entities.User
		repository_res int
		err            error
	}{
		{
			test_name: "user created successfully",
			user: entities.User{
				Email:    "user@email.com",
				Password: "qwerty",
				Age:      23,
				Details:  repository.GenereateDetails(),
			},
			repository_res: 1,
		},
		{
			test_name: "no password error",
			user: entities.User{
				Email:   "user@email.com",
				Age:     23,
				Details: repository.GenereateDetails(),
			},
			repository_res: -1,
		},
		{
			test_name: "no email error",
			user: entities.User{
				Password: "qwerty",
				Age:      23,
				Details:  repository.GenereateDetails(),
			},
			repository_res: -1,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			//  prepare
			assert := assert.New(t)
			ctx := context.Background()

			user_req := &userpb.CreateUserRequest{
				Email:    tc.user.Email,
				Password: tc.user.Password,
				Age:      uint32(tc.user.Age),
			}

			details_req := &detailspb.SetUserDetailsRequest{
				UserId:       uint32(tc.repository_res),
				Country:      tc.user.Country,
				City:         tc.user.City,
				MobileNumber: tc.user.MobileNumber,
				Married:      tc.user.Married,
				Height:       tc.user.Height,
				Weight:       tc.user.Weight,
			}

			user_res := &userpb.CreateUserResponse{Id: int32(tc.repository_res)}
			details_res := &detailspb.SetUserDetailsResponse{Success: tc.repository_res == 1}

			user_mock.On("CreateUser", mock.Anything, user_req).Return(user_res, tc.err)
			details_mock.On("SetUserDetails", mock.Anything, details_req).Return(details_res, tc.err)

			// act
			res, err := http_repository.CreateUser(ctx, tc.user)

			// assert
			assert.Equal(res, tc.repository_res)
			assert.Equal(err, tc.err)
		})
	}
}

func TestAuthenticate(t *testing.T) {
	user_mock := new(repository.GrpcUserMock)
	details_mock := new(repository.GrpcDetailsMock)
	conn1, conn2, _ := repository.InitRepoMock(user_mock, details_mock)

	defer conn1.Close()
	defer conn2.Close()

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
			user_mock.On("Authenticate", ctx, grpc_req).Return(grpc_res, tc.err)

			// act
			res, err := user_mock.Authenticate(ctx, grpc_req)

			// assert
			assert.Equal(res.GetUserId(), int32(tc.res))
			assert.Equal(err, tc.err)
		})
	}
}

func UpdateUser(t *testing.T) {
	user_mock := new(repository.GrpcUserMock)
	details_mock := new(repository.GrpcDetailsMock)
	conn1, conn2, http_repository := repository.InitRepoMock(user_mock, details_mock)

	defer conn1.Close()
	defer conn2.Close()

	test_cases := []struct {
		test_name string
		data      entities.UserUpdate
		res       bool
		err       error
	}{
		{
			test_name: "Update email successfully",
			data: entities.UserUpdate{
				UserId: 0,
				User: entities.User{
					Email:    "email@domian.com",
					Password: "qwerty",
					Age:      23,
					Details:  repository.GenereateDetails(),
				},
			},
			res: true,
			err: nil,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			//  prepare
			ctx := context.Background()
			assert := assert.New(t)

			user_req := &userpb.UpdateUserRequest{
				Id:       uint32(tc.data.Age),
				Email:    tc.data.Email,
				Password: tc.data.Password,
				Age:      uint32(tc.data.Age),
			}
			details_req := &detailspb.SetUserDetailsRequest{
				UserId:       uint32(tc.data.UserId),
				Country:      tc.data.Country,
				City:         tc.data.City,
				MobileNumber: tc.data.MobileNumber,
				Married:      tc.data.Married,
				Height:       tc.data.Height,
				Weight:       tc.data.Weight,
			}

			user_res := &userpb.UpdateUserResponse{Success: tc.res}
			details_res := &detailspb.SetUserDetailsResponse{Success: tc.res}

			user_mock.On("UpdateUser", ctx, user_req).Return(user_res, tc.err)
			details_mock.On("SetUserDetails", mock.Anything, details_req).Return(details_res, tc.err)

			// act
			res, err := http_repository.UpdateUser(ctx, tc.data)

			// assert
			assert.Equal(res, user_res)
			assert.Equal(err, tc.err)
		})
	}
}

func TestGetUser(t *testing.T) {
	user_mock := new(repository.GrpcUserMock)
	details_mock := new(repository.GrpcDetailsMock)
	conn1, conn2, _ := repository.InitRepoMock(user_mock, details_mock)

	defer conn1.Close()
	defer conn2.Close()

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
				Email:    "email@domain.com",
				Password: "password",
				Age:      10,
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
				Email:    tc.res.Email,
				Password: tc.res.Password,
				Age:      uint32(tc.res.Age),
			}
			user_mock.On("GetUser", ctx, grpc_req).Return(grpc_res, tc.err)

			// act
			res, err := user_mock.GetUser(ctx, grpc_req)

			// assert
			assert.Equal(res, grpc_res)
			assert.Equal(err, tc.err)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	user_mock := new(repository.GrpcUserMock)
	details_mock := new(repository.GrpcDetailsMock)
	conn1, conn2, _ := repository.InitRepoMock(user_mock, details_mock)

	defer conn1.Close()
	defer conn2.Close()

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

			user_mock.On("DeleteUser", ctx, grpc_req).Return(grpc_res, tc.err)

			// act
			res, err := user_mock.DeleteUser(ctx, grpc_req)

			// assert
			assert.Equal(res, grpc_res)
			assert.Equal(err, tc.err)
		})
	}
}
