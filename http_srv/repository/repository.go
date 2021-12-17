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
	UpdateUser(ctx context.Context, user entities.UserUpdate) (entities.User, error)
	GetUser(ctx context.Context, id int) (entities.User, error)
	DeleteUser(ctx context.Context, id int) (bool, error)
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

	grpc_req := userpb.CreateUserRequest{
		Email:                 user.Email,
		Password:              user.Password,
		Age:                   uint32(user.Age),
		AdditionalInformation: user.ExtraInfo,
	}
	grpc_res, err := hr.client.CreateUser(ctx, &grpc_req)

	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	return grpc_res.GetId(), nil
}

func (hr httpRepository) Authenticate(ctx context.Context, session entities.Session) (int, error) {
	logger := log.With(hr.logger, "method", "authenticate_user")

	if session.Email == "" || session.Password == "" {
		return -1, errors.New("Email or Password empty")
	}

	grpc_req := userpb.AuthenticateRequest{
		Email:    session.Email,
		Password: session.Password,
	}
	grpc_res, err := hr.client.Authenticate(ctx, &grpc_req)

	if err != nil {
		level.Error(logger).Log("err", err)
		return -1, err
	}

	return int(grpc_res.GetUserId()), nil
}

func (hr httpRepository) UpdateUser(ctx context.Context, user entities.UserUpdate) (entities.User, error) {
	logger := log.With(hr.logger, "method", "update_user")

	grpc_req := userpb.UpdateUserRequest{
		Id:                    uint32(user.UserId),
		Email:                 user.Email,
		Password:              user.Password,
		Age:                   uint32(user.Age),
		AdditionalInformation: user.ExtraInfo,
	}

	grpc_res, err := hr.client.UpdateUser(ctx, &grpc_req)

	if err != nil {
		level.Error(logger).Log("err", err)
		return entities.User{}, err
	}

	u := entities.User{
		Email:     grpc_res.GetEmail(),
		Password:  grpc_res.GetPassword(),
		Age:       int(grpc_res.GetAge()),
		ExtraInfo: grpc_res.AdditionalInformation,
	}

	return u, nil
}

func (hr httpRepository) GetUser(ctx context.Context, id int) (entities.User, error) {
	logger := log.With(hr.logger, "method", "update_user")

	grpc_req := userpb.GetUserRequest{
		Id: uint32(id),
	}

	grpc_res, err := hr.client.GetUser(ctx, &grpc_req)

	if err != nil {
		level.Error(logger).Log("err", err)
		return entities.User{}, err
	}

	u := entities.User{
		Email:     grpc_res.GetEmail(),
		Password:  grpc_res.GetPassword(),
		Age:       int(grpc_res.GetAge()),
		ExtraInfo: grpc_res.GetAdditionalInformation(),
	}

	return u, nil
}

func (hr httpRepository) DeleteUser(ctx context.Context, id int) (bool, error) {
	logger := log.With(hr.logger, "method", "delete_user")

	grpc_req := userpb.DeleteUserRequest{
		Id: uint32(id),
	}

	grpc_res, err := hr.client.DeleteUser(ctx, &grpc_req)

	if err != nil {
		level.Error(logger).Log("err", err)
		return false, err
	}

	return grpc_res.GetSuccess(), nil
}
