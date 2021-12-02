package service

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/http_srv/repository"
)

var (
	Empty_Field = errors.New("Email or Password empy!")
)

type HttpService interface {
	CreateUser(ctx context.Context, email string, pwd string, extra_info string, age int) (string, error)
}

type httpService struct {
	logger     log.Logger
	repository repository.HttpRepository
}

func NewHttpService(r repository.HttpRepository, l log.Logger) HttpService {
	return &httpService{
		logger:     l,
		repository: r,
	}
}

func (hs httpService) CreateUser(ctx context.Context, email string, pwd string, extra_info string, age int) (string, error) {
	logger := log.With(hs.logger, "HTTP_SRV: method", "create_user")

	if email == "" || pwd == "" {
		return "", Empty_Field
	}

	res, err := hs.repository.CreateUser(ctx, email, pwd, extra_info, age)

	if err != nil {
		level.Error(logger).Log("ERROR: ", err)
		return "", err
	}

	logger.Log("user_created_successfully", res)
	return res, nil
}
