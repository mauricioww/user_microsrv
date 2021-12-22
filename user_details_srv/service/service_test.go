package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/mauricioww/user_microsrv/user_details_srv/entities"
	"github.com/mauricioww/user_microsrv/user_details_srv/service"
	"github.com/stretchr/testify/assert"
)

func TestSetUserDetails(t *testing.T) {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(
			logger,
			"service",
			"user_details",
			"time",
			log.DefaultTimestampUTC,
			"caller",
			log.DefaultCaller,
		)
	}

	var grpc_user_details_srv service.GrpcUserDetailsService

	user_details_repo_mock := new(service.UserDetailsRepositoryMock)
	grpc_user_details_srv = service.NewGrpcUserDetailsService(user_details_repo_mock, logger)

	test_cases := []struct {
		test_name string
		data      entities.UserDetails
		res       bool
		err       error
	}{
		{
			test_name: "set user details which no exists success",
			data: entities.UserDetails{
				UserId:       1,
				Country:      "Mexico",
				City:         "CDMX",
				MobileNumber: "11223344",
				Married:      false,
				Height:       1.75,
				Weigth:       76.0,
			},
			res: true,
			err: nil,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			user_details_repo_mock.On("SetUserDetails", ctx, tc.data).Return(tc.res, tc.err)
			res, err := grpc_user_details_srv.SetUserDetails(ctx, tc.data.UserId, tc.data.Country, tc.data.City,
				tc.data.MobileNumber, tc.data.Married, tc.data.Height, tc.data.Weigth)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}