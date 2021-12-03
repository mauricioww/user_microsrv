package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/mauricioww/user_microsrv/user_srv/repository"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
)

type GrpcUserService interface {
	CreateUser(ctx context.Context, user *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error)
}

type grpcUserService struct {
	reepository repository.UserSrvRepository
	logger      log.Logger
}

func NewGrpcUserService(l log.Logger, r repository.UserSrvRepository) GrpcUserService {
	return &grpcUserService{
		logger:      l,
		reepository: r,
	}
}

func (g grpcUserService) CreateUser(ctx context.Context, user_pb *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	logger := log.With(g.logger, "GRPC_USER_SERVICE: methd", "create_user")

	user := entities.User{
		Email:     user_pb.GetEmail(),
		Password:  user_pb.GetPassword(),
		Age:       int(user_pb.GetAge()),
		ExtraInfo: user_pb.GetAdditionalInformation(),
	}
	res, err := g.reepository.CreateUser(ctx, user)

	if err != nil {
		level.Error(logger).Log("err", err)
		return &userpb.CreateUserResponse{}, err
	}

	logger.Log("create user", res)
	return &userpb.CreateUserResponse{Id: res}, nil
}
