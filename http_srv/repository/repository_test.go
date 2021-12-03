package repository_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/http_srv/repository"
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

	http_repository := repository.NewHttpRepository(conn, logger)

	test_cases := []struct {
		test_name string
		user      entities.User
		res       string
		err       error
	}{
		{
			test_name: "user created successfully",
			user: entities.User{
				Email:     "user@email.com",
				Password:  "qwerty",
				Age:       23,
				ExtraInfo: "fav movie: fight club",
			},
			res: "success",
			err: nil,
		},
		{
			test_name: "no password",
			user: entities.User{
				Email:     "user@email.com",
				Age:       23,
				ExtraInfo: "fav movie: fight club",
			},
			res: "",
			err: errors.New("Email or Password empty!"),
		},
		{
			test_name: "no email",
			user: entities.User{
				Password:  "qwerty",
				Age:       23,
				ExtraInfo: "fav movie: fight club",
			},
			res: "",
			err: errors.New("Email or Password empty!"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			//  prepare
			ctx := context.Background()
			assert := assert.New(t)
			grpc_mock.On("CreateUser", ctx, tc.user).Return(tc.res, tc.err)

			// act
			res, err := http_repository.CreateUser(ctx, tc.user)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(err, tc.err)
		})
	}
}
