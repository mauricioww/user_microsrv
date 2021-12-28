package transport_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mauricioww/user_microsrv/user_details_srv/entities"
	"github.com/mauricioww/user_microsrv/user_details_srv/transport"
	"github.com/stretchr/testify/assert"
)

func TestSetUserDetails(t *testing.T) {
	grpc_user_srv_mock := new(transport.GrpcUserDetailsSrvMock)

	endpoints := transport.MakeGrpcUserDetailsServiceEndpoints(grpc_user_srv_mock)

	test_cases := []struct {
		test_name string
		data      transport.SetUserDetailsRequest
		res       transport.SetUserDetailsResponse
		srv_res   bool
		err       error
	}{
		{
			test_name: "set user details first time values success",
			data: transport.SetUserDetailsRequest{
				UserId:       1,
				Country:      "Mexico",
				City:         "CDMX",
				MobileNumber: "11223344",
				Married:      false,
				Height:       1.75,
				Weigth:       76.0,
			},
			res: transport.SetUserDetailsResponse{
				Success: true,
			},
			srv_res: true,
			err:     nil,
		},
		{
			test_name: "set user details specific values success",
			data: transport.SetUserDetailsRequest{
				UserId:       1,
				MobileNumber: "12345789",
			},
			res: transport.SetUserDetailsResponse{
				Success: true,
			},
			srv_res: true,
			err:     nil,
		},
	}
	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()

			// act
			grpc_user_srv_mock.On("SetUserDetails", ctx, tc.data.UserId, tc.data.Country, tc.data.City,
				tc.data.MobileNumber, tc.data.Married, tc.data.Height, tc.data.Weigth).Return(tc.srv_res, tc.err)

			res, err := endpoints.SetUserDetails(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestGetUserDetails(t *testing.T) {
	grpc_user_srv_mock := new(transport.GrpcUserDetailsSrvMock)

	endpoints := transport.MakeGrpcUserDetailsServiceEndpoints(grpc_user_srv_mock)

	test_cases := []struct {
		test_name string
		data      transport.GetUserDetailsRequest
		res       transport.GetUserDetailsResponse
		srv_res   entities.UserDetails
		err       error
	}{
		{
			test_name: "get user details success",
			data: transport.GetUserDetailsRequest{
				UserId: 0,
			},
			srv_res: entities.UserDetails{
				Country:      "Mexico",
				City:         "CDMX",
				MobileNumber: "11223344",
				Married:      false,
				Height:       1.75,
				Weight:       76.0,
			},
			err: nil,
		},
		{
			test_name: "get user details which does not exist error",
			data: transport.GetUserDetailsRequest{
				UserId: 1,
			},
			err: errors.New("User not found"),
		},
	}
	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()
			tc.res = transport.GetUserDetailsResponse{Country: tc.srv_res.Country, City: tc.srv_res.City, MobileNumber: tc.srv_res.MobileNumber,
				Married: tc.srv_res.Married, Height: tc.srv_res.Height, Weight: tc.srv_res.Weight}

			// act
			grpc_user_srv_mock.On("GetUserDetails", ctx, tc.data.UserId).Return(tc.srv_res, tc.err)

			res, err := endpoints.GetUserDetails(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestDeleteUserDetails(t *testing.T) {
	grpc_user_srv_mock := new(transport.GrpcUserDetailsSrvMock)

	endpoints := transport.MakeGrpcUserDetailsServiceEndpoints(grpc_user_srv_mock)

	test_cases := []struct {
		test_name string
		data      transport.DeleteUserDetailsRequest
		res       transport.DeleteUserDetailsResponse
		srv_res   bool
		err       error
	}{
		{
			test_name: "delete user details success",
			data: transport.DeleteUserDetailsRequest{
				UserId: 0,
			},
			srv_res: true,
			err:     nil,
		},
		{
			test_name: "delete user details which does not exist error",
			data: transport.DeleteUserDetailsRequest{
				UserId: 1,
			},
			err: errors.New("User not found"),
		},
	}
	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()
			tc.res = transport.DeleteUserDetailsResponse{Success: tc.srv_res}

			// act
			grpc_user_srv_mock.On("DeleteUserDetails", ctx, tc.data.UserId).Return(tc.srv_res, tc.err)

			res, err := endpoints.DeleteUserDetails(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}
