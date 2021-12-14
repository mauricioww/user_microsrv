package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mauricioww/user_microsrv/http_srv/service"
)

type HttpEndpoints struct {
	CreateUser   endpoint.Endpoint
	Authenticate endpoint.Endpoint
	UpdateUser   endpoint.Endpoint
}

func MakeHttpEndpoints(http_srv service.HttpService) HttpEndpoints {
	return HttpEndpoints{
		CreateUser:   makeCreateUserEndpoint(http_srv),
		Authenticate: makeAuthenticateEndpoint(http_srv),
		UpdateUser:   makeUpdateUserEndpoint(http_srv),
	}
}

func makeCreateUserEndpoint(http_srv service.HttpService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		res, err := http_srv.CreateUser(ctx, req.Email, req.Password, req.ExtraInfo, req.Age)
		return CreateUserResponse{Id: res, Email: req.Email, Password: req.Password, Age: req.Age, ExtraInfo: req.ExtraInfo}, err
	}
}

func makeAuthenticateEndpoint(http_srv service.HttpService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthenticateRequest)
		res, err := http_srv.Authenticate(ctx, req.Email, req.Password)
		return AuthenticateResponse{Token: res}, err
	}
}

func makeUpdateUserEndpoint(http_srv service.HttpService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateUserRequest)
		res, err := http_srv.UpdateUser(ctx, req.UserId, req.Email, req.Password, req.ExtraInfo, req.Age)
		return UpdateUserResponse{Id: req.UserId, Email: res.Email, Password: req.Password, Age: res.Age, ExtraInfo: res.ExtraInfo}, err
	}
}
