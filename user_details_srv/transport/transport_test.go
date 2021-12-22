package transport_test

import (
	"context"
	"testing"

	"github.com/mauricioww/user_microsrv/user_details_srv/transport"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
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
