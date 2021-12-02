package repository

import (
	"context"

	"github.com/go-kit/kit/log"
)

type HttpRepository interface {
	CreateUser(ctx context.Context, email string, pwd string, extra_info string, age int) (string, error)
}

type httpRepository struct {
	logger log.Logger
}

func NewHttpRepository(logger log.Logger) HttpRepository {
	return httpRepository{
		logger: log.With(logger, "http_repository", "proxy?"),
	}
}

func (hr httpRepository) CreateUser(ctx context.Context, email string, pwd string, extra_info string, age int) (string, error) {
	return email, nil
}
