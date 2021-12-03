package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mauricioww/user_microsrv/http_srv/service"
)

type HttpEndpoints struct {
	CreateUser endpoint.Endpoint
}

func MakeHttpEndpoints(http_srv service.HttpService) HttpEndpoints {
	return HttpEndpoints{
		CreateUser: makeCreateUserEndpoint(http_srv),
	}
}

func makeCreateUserEndpoint(http_srv service.HttpService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		res, err := http_srv.CreateUser(ctx, req.Email, req.Password, req.ExtraInfo, req.Age)
		return CreateUserResponse{ID: res}, err
	}
}
