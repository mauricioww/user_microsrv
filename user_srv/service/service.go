package service

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/helpers"
	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/mauricioww/user_microsrv/user_srv/repository"
)

type GrpcUserService interface {
	CreateUser(ctx context.Context, email string, pwd string, extra_info string, age int) (string, error)
	Authenticate(ctx context.Context, email string, pwd string) (int, error)
	UpdateUser(ctx context.Context, id int, email string, pwd string, extra_info string, age int) (entities.User, error)
	GetUser(ctx context.Context, id int) (entities.User, error)
	DeleteUser(ctx context.Context, id int) (bool, error)
}

type grpcUserService struct {
	repository repository.UserSrvRepository
	logger     log.Logger
}

func NewGrpcUserService(r repository.UserSrvRepository, l log.Logger) GrpcUserService {
	return &grpcUserService{
		logger:     l,
		repository: r,
	}
}

func (g *grpcUserService) CreateUser(ctx context.Context, email string, pwd string, extra_info string, age int) (string, error) {
	logger := log.With(g.logger, "method", "create_user")
	ciphered_pwd := helpers.Cipher(pwd)

	user := entities.User{
		Email:     email,
		Password:  ciphered_pwd,
		Age:       age,
		ExtraInfo: extra_info,
	}

	res, err := g.repository.CreateUser(ctx, user)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
		return "", err
	}

	logger.Log("action", "success")
	return res, nil
}

func (g *grpcUserService) Authenticate(ctx context.Context, email string, pwd string) (int, error) {
	logger := log.With(g.logger, "method", "authenticate")

	auth := entities.Session{
		Email:    email,
		Password: pwd,
	}

	hashed_pwd, err := g.repository.Authenticate(ctx, &auth)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
		return -1, err
	}

	if cipher := helpers.Cipher(auth.Password); cipher != hashed_pwd {
		return -1, errors.New("Password error")
	}

	logger.Log("action", "success")
	return auth.Id, nil
}

func (g *grpcUserService) UpdateUser(ctx context.Context, id int, email string, pwd string, extra_info string, age int) (entities.User, error) {
	logger := log.With(g.logger, "method", "update_user")
	ciphered_pwd := helpers.Cipher(pwd)

	update_info := entities.Update{
		UserId: id,
		User: entities.User{
			Email:     email,
			Password:  ciphered_pwd,
			ExtraInfo: extra_info,
			Age:       age,
		},
	}

	user, err := g.repository.UpdateUser(ctx, update_info)
	user.Password = pwd
	if err != nil {
		level.Error(logger).Log("ERROR", err)
		return entities.User{}, err
	}

	logger.Log("action", "success")
	return user, nil
}

func (g *grpcUserService) GetUser(ctx context.Context, id int) (entities.User, error) {
	logger := log.With(g.logger, "method", "get_user")

	user, err := g.repository.GetUser(ctx, id)
	if err != nil {
		level.Error(logger).Log("ERROR", err)
		return entities.User{}, err
	}

	original_pwd := helpers.Decipher(user.Password)
	user.Password = original_pwd

	logger.Log("action", "success")
	return user, nil
}

func (g *grpcUserService) DeleteUser(ctx context.Context, id int) (bool, error) {
	logger := log.With(g.logger, "method", "delete_user")

	success, err := g.repository.DeleteUser(ctx, id)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
	}

	logger.Log("action", "success")
	return success, err
}
