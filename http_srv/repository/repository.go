package repository

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
	"google.golang.org/grpc"
)

type HttpRepository interface {
	CreateUser(ctx context.Context, user entities.User) (string, error)
	Authenticate(ctx context.Context, session entities.Session) (int, error)
}

type httpRepository struct {
	client userpb.UserServiceClient
	logger log.Logger
}

func NewHttpRepository(conn *grpc.ClientConn, logger log.Logger) HttpRepository {
	return httpRepository{
		client: userpb.NewUserServiceClient(conn),
		logger: log.With(logger, "http_repository", "proxy?"),
	}
}

func (hr httpRepository) CreateUser(ctx context.Context, user entities.User) (string, error) {
	logger := log.With(hr.logger, "method", "create_users")

	if user.Email == "" || user.Password == "" {
		return "", errors.New("Email or Password empty!")
	}

	grpc_request := userpb.CreateUserRequest{
		Email:                 user.Email,
		Password:              user.Password,
		Age:                   uint32(user.Age),
		AdditionalInformation: user.ExtraInfo,
	}
	grpc_response, err := hr.client.CreateUser(ctx, &grpc_request)

	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	res := grpc_response.GetId()
	return res, nil
}

func (hr httpRepository) Authenticate(ctx context.Context, session entities.Session) (int, error) {
	logger := log.With(hr.logger, "method", "create_users")

	if session.Email == "" || session.Password == "" {
		return -1, errors.New("Email or Password empty")
	}

	grpc_request := userpb.AuthenticateRequest{
		Email:    session.Email,
		Password: session.Password,
	}
	grpc_response, err := hr.client.Authenticate(ctx, &grpc_request)

	if err != nil {
		level.Error(logger).Log("err", err)
		return -1, err
	}

	return int(grpc_response.GetUserId()), nil
}
