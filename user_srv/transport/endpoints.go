package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mauricioww/user_microsrv/user_srv/service"
)

type GrpcUserServiceEndpoints struct {
	CreateUser endpoint.Endpoint
}

func MakeGrpcUserServiceEndpoints(grpc_user_srv service.GrpcUserService) GrpcUserServiceEndpoints {
	return GrpcUserServiceEndpoints{
		CreateUser: makeCreateUserEndpoint(grpc_user_srv),
	}
}

func makeCreateUserEndpoint(srv service.GrpcUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, _ := request.(CreateUserRequest)
		res, err := srv.CreateUser(ctx, req.Email, req.Password, req.ExtraInfo, req.Age)
		return res, err
	}
}
