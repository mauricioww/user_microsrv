package service

import (
	"context"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/http_srv/repository"
)

type HttpService interface {
	CreateUser(ctx context.Context, email string, pwd string, age int, details entities.Details) (int, error)
	Authenticate(ctx context.Context, email string, pwd string) (string, error)
	UpdateUser(ctx context.Context, user_id int, email string, pwd string, age int, details entities.Details) (entities.User, error)
	GetUser(ctx context.Context, user_id int) (entities.User, error)
	DeleteUser(ctx context.Context, user_id int) (bool, error)
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

func (hs httpService) CreateUser(ctx context.Context, email string, pwd string, age int, details entities.Details) (int, error) {
	logger := log.With(hs.logger, "method", "create_user")

	user := entities.User{
		Email:    email,
		Password: pwd,
		Age:      age,
		Details:  details,
	}

	res, err := hs.repository.CreateUser(ctx, user)

	if err != nil {
		level.Error(logger).Log("ERROR: ", err)
		return -1, err
	}

	logger.Log("action", "success")
	return res, nil
}

func (hs httpService) Authenticate(ctx context.Context, email string, pwd string) (string, error) {
	logger := log.With(hs.logger, "method", "authenticate")

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

	if res >= 0 {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": strconv.Itoa(res),
			"email":   session.Email,
			"exp":     time.Now().Add(time.Minute * 15).Unix(),
		})

		response, err = token.SignedString([]byte("this_is_a_secret_shhh"))

		if err != nil {
			level.Error(logger).Log("ERROR: ", err)
			return "", err
		}
	}

	logger.Log("action", "success")
	return response, nil
}

func (hs httpService) UpdateUser(ctx context.Context, user_id int, email string, pwd string, age int, details entities.Details) (entities.User, error) {
	logger := log.With(hs.logger, "method", "authenticate")
	var res entities.User

	info_update := entities.UserUpdate{
		UserId: user_id,
		User: entities.User{
			Email:    email,
			Password: pwd,
			Age:      age,
			Details:  details,
		},
	}

	success, err := hs.repository.UpdateUser(ctx, info_update)

	if err != nil {
		level.Error(logger).Log("ERROR: ", err)
		return entities.User{}, err
	}

	if success {
		res = info_update.User
	}

	logger.Log("action", "success")
	return res, err
}

func (hs httpService) GetUser(ctx context.Context, user_id int) (entities.User, error) {
	logger := log.With(hs.logger, "method", "get_user")

	res, err := hs.repository.GetUser(ctx, user_id)

	if err != nil {
		level.Error(logger).Log("ERROR: ", err)
		return entities.User{}, err
	}

	logger.Log("action", "success")
	return res, nil
}

func (hs httpService) DeleteUser(ctx context.Context, user_id int) (bool, error) {
	logger := log.With(hs.logger, "method", "delete_user")

	res, err := hs.repository.DeleteUser(ctx, user_id)

	if err != nil {
		level.Error(logger).Log("ERROR: ", err)
		return false, err
	}

	logger.Log("action", "success")
	return res, nil
}
