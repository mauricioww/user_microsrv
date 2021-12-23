package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/user_details_srv/detailspb"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
	"google.golang.org/grpc"
)

type HttpRepository interface {
	CreateUser(ctx context.Context, user entities.User) (int, error)
	Authenticate(ctx context.Context, session entities.Session) (int, error)
	UpdateUser(ctx context.Context, user entities.UserUpdate) (bool, error)
	GetUser(ctx context.Context, id int) (entities.User, error)
	DeleteUser(ctx context.Context, id int) (bool, error)
}

type httpRepository struct {
	user_client    userpb.UserServiceClient
	details_client detailspb.UserDetailsServiceClient
	logger         log.Logger
}

func NewHttpRepository(conn1 *grpc.ClientConn, conn2 *grpc.ClientConn, logger log.Logger) HttpRepository {
	return httpRepository{
		user_client:    userpb.NewUserServiceClient(conn1),
		details_client: detailspb.NewUserDetailsServiceClient(conn2),
		logger:         log.With(logger, "http_repository", "proxy?"),
	}
}

func (hr httpRepository) CreateUser(ctx context.Context, user entities.User) (int, error) {
	logger := log.With(hr.logger, "method", "create_users")

	userpb_req := userpb.CreateUserRequest{
		Email:    user.Email,
		Password: user.Password,
		Age:      uint32(user.Age),
	}

	user_res, err := hr.user_client.CreateUser(ctx, &userpb_req)
	if err != nil {
		level.Error(logger).Log("err", err)
		return -1, err
	}

	detailspb_req := detailspb.SetUserDetailsRequest{
		UserId:       uint32(user_res.GetId()),
		Country:      user.Details.Country,
		City:         user.Details.City,
		MobileNumber: user.Details.MobileNumber,
		Married:      user.Details.Married,
		Height:       user.Details.Height,
		Weigth:       user.Details.Weigth,
	}

	details_res, err := hr.details_client.SetUserDetails(ctx, &detailspb_req)
	if err != nil {
		level.Error(logger).Log("err", err)
		return -1, err
	}

	if details_res.GetSuccess() {
		fmt.Println("details setted")
	}

	return int(user_res.Id), nil
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
	grpc_res, err := hr.user_client.Authenticate(ctx, &grpc_req)

	if err != nil {
		level.Error(logger).Log("err", err)
		return -1, err
	}

	return int(grpc_res.GetUserId()), nil
}

func (hr httpRepository) UpdateUser(ctx context.Context, user entities.UserUpdate) (bool, error) {
	logger := log.With(hr.logger, "method", "update_user")
	var res bool

	user_req := userpb.UpdateUserRequest{
		Id:       uint32(user.UserId),
		Email:    user.Email,
		Password: user.Password,
		Age:      uint32(user.Age),
	}
	details_req := detailspb.SetUserDetailsRequest{
		UserId:       uint32(user.UserId),
		Country:      user.Details.Country,
		City:         user.Details.City,
		MobileNumber: user.Details.MobileNumber,
		Married:      user.Details.Married,
		Height:       user.Details.Height,
		Weigth:       user.Details.Weigth,
	}

	user_res, err := hr.user_client.UpdateUser(ctx, &user_req)
	if err != nil {
		level.Error(logger).Log("err", err)
		return false, err
	}

	res = user_res.GetSuccess()

	details_res, err := hr.details_client.SetUserDetails(ctx, &details_req)
	if err != nil {
		level.Error(logger).Log("err", err)
		return false, err
	}

	res = details_res.GetSuccess()

	return res, nil
}

func (hr httpRepository) GetUser(ctx context.Context, id int) (entities.User, error) {
	logger := log.With(hr.logger, "method", "update_user")

	grpc_req := userpb.GetUserRequest{
		Id: uint32(id),
	}

	grpc_res, err := hr.user_client.GetUser(ctx, &grpc_req)

	if err != nil {
		level.Error(logger).Log("err", err)
		return entities.User{}, err
	}

	u := entities.User{
		Email:    grpc_res.GetEmail(),
		Password: grpc_res.GetPassword(),
		Age:      int(grpc_res.GetAge()),
	}

	return u, nil
}

func (hr httpRepository) DeleteUser(ctx context.Context, id int) (bool, error) {
	logger := log.With(hr.logger, "method", "delete_user")

	grpc_req := userpb.DeleteUserRequest{
		Id: uint32(id),
	}

	grpc_res, err := hr.user_client.DeleteUser(ctx, &grpc_req)

	if err != nil {
		level.Error(logger).Log("err", err)
		return false, err
	}

	return grpc_res.GetSuccess(), nil
}
