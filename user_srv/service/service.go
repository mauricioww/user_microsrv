package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/mauricioww/user_microsrv/user_srv/repository"
)

type GrpcUserService interface {
	CreateUser(ctx context.Context, email string, pwd string, extra_info string, age int) (string, error)
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

func (g *grpcUserService) CreateUser(ctx context.Context, email string, pwd string, extra_info string, age int) (string, error) {
	logger := log.With(g.logger, "GRPC_USER_SERVICE: method", "create_user")

	user := entities.User{
		Email:     email,
		Password:  pwd,
		Age:       age,
		ExtraInfo: extra_info,
	}

	res, err := g.reepository.CreateUser(ctx, user)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
		return "", err
	}

	logger.Log("user_saved_successfully", res)
	return res, nil
}