package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/http_srv/repository"
)

type HttpService interface {
	CreateUser(ctx context.Context, email string, pwd string, extra_info string, age int) (string, error)
}

type httpService struct {
	repository repository.HttpRepository
	logger     log.Logger
}

func NewHttpService(r repository.HttpRepository, l log.Logger) HttpService {
	return &httpService{
		logger:     l,
		repository: r,
	}
}

func (hs httpService) CreateUser(ctx context.Context, email string, pwd string, extra_info string, age int) (string, error) {
	logger := log.With(hs.logger, "HTTP_SRV: method", "create_user")
	// hashed_pwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)

	user := entities.User{
		Email:     email,
		Password:  pwd,
		Age:       age,
		ExtraInfo: extra_info,
	}

	res, err := hs.repository.CreateUser(ctx, user)

	if err != nil {
		level.Error(logger).Log("ERROR: ", err)
		return "", err
	}

	logger.Log("user_send_successfully", res)
	return res, nil
}
