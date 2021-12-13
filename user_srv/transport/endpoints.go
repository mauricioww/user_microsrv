package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mauricioww/user_microsrv/user_srv/service"
)

type GrpcUserServiceEndpoints struct {
	CreateUser   endpoint.Endpoint
	Authenticate endpoint.Endpoint
	UpdateUser   endpoint.Endpoint
}

func MakeGrpcUserServiceEndpoints(grpc_user_srv service.GrpcUserService) GrpcUserServiceEndpoints {
	return GrpcUserServiceEndpoints{
		CreateUser:   makeCreateUserEndpoint(grpc_user_srv),
		Authenticate: makeAuthenticateEndpoint(grpc_user_srv),
		UpdateUser:   makeUpdateUserEndpoint(grpc_user_srv),
	}
}

func makeCreateUserEndpoint(srv service.GrpcUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, _ := request.(CreateUserRequest)
		res, err := srv.CreateUser(ctx, req.Email, req.Password, req.ExtraInfo, req.Age)
		return res, err
	}
}

func makeAuthenticateEndpoint(srv service.GrpcUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, _ := request.(AuthenticateRequest)
		res, err := srv.Authenticate(ctx, req.Email, req.Password)
		return res, err
	}
}

func makeUpdateUserEndpoint(srv service.GrpcUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, _ := request.(UpdateUserRequest)
		res, err := srv.UpdateUser(ctx, req.Id, req.Information)
		return res, err
	}
}
