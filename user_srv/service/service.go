package service

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/mauricioww/user_microsrv/user_srv/repository"
	"golang.org/x/crypto/bcrypt"
)

type GrpcUserService interface {
	CreateUser(ctx context.Context, email string, pwd string, extra_info string, age int) (string, error)
	Authenticate(ctx context.Context, email string, pwd string) (string, error)
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
	logger := log.With(g.logger, "GRPC_USER_SERVICE: method", "create_user")
	hashed_pwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)

	user := entities.User{
		Email:     email,
		Password:  string(hashed_pwd),
		Age:       age,
		ExtraInfo: extra_info,
	}

	res, err := g.repository.CreateUser(ctx, user)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
		return "", err
	}

	logger.Log("user_saved_successfully", res)
	return res, nil
}

func (g *grpcUserService) Authenticate(ctx context.Context, email string, pwd string) (string, error) {
	logger := log.With(g.logger, "GRPC_USER_SERVICE: method", "authenticate")

	auth := entities.Session{
		Email:    email,
		Password: pwd,
	}

	hashed_pwd, err := g.repository.Authenticate(ctx, auth)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashed_pwd), []byte(auth.Password)); err != nil {
		return "", errors.New("Password error")
	}

	return "user_authenticated", nil
}
