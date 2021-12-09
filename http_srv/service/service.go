package service

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/http_srv/repository"
)

type HttpService interface {
	CreateUser(ctx context.Context, email string, pwd string, extra_info string, age int) (string, error)
	Authenticate(ctx context.Context, email string, pwd string) (string, error)
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

func (hs httpService) Authenticate(ctx context.Context, email string, pwd string) (string, error) {
	logger := log.With(hs.logger, "HTTP_SRV: method", "authenticate")
	var response string
	session := entities.Session{
		Email:    email,
		Password: pwd,
	}

	res, err := hs.repository.Authenticate(ctx, session)

	if err != nil {
		level.Error(logger).Log("ERROR: ", err)
		return "", err
	}

	if res == "user_authenticated" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Minute * 15).Unix(),
		})

		response, err = token.SignedString([]byte("this_is_a_secret_shhh"))

		if err != nil {
			level.Error(logger).Log("ERROR: ", err)
			return "", err
		}

		logger.Log("user authenticated", res)
	}

	return response, nil
}
