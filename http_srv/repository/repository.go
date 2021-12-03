package repository

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
	"google.golang.org/grpc"
)

type HttpRepository interface {
	CreateUser(ctx context.Context, user entities.User) (string, error)
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
	if user.Email == "" || user.Password == "" {
		return "", errors.New("Email or Password empty!")
	}

	return "success", nil
}
