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
	Authenticate(ctx context.Context, session entities.Session) (string, error)
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

	request := userpb.CreateUserRequest{
		Email:                 user.Email,
		Password:              user.Password,
		Age:                   uint32(user.Age),
		AdditionalInformation: user.ExtraInfo,
	}
	grpc_response, err := hr.client.CreateUser(ctx, &request)

	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	res := grpc_response.GetId()
	return res, nil
}

func (hr httpRepository) Authenticate(ctx context.Context, session entities.Session) (string, error) {
	logger := log.With(hr.logger, "method", "create_users")

	if session.Email == "" || session.Password == "" {
		return "", errors.New("Email or Password empty!")
	}

	// TODO: send grpc method
	var err error = nil

	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	// TODO: return a real token
	return "auth_token", nil
}
