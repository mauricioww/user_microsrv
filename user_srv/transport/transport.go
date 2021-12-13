package transport

import (
	"context"
	"errors"

	grpc_gokit "github.com/go-kit/kit/transport/grpc"
	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
)

type gRPCServer struct {
	createUser   grpc_gokit.Handler
	authenticate grpc_gokit.Handler
	updateUser   grpc_gokit.Handler
	userpb.UnimplementedUserServiceServer
}

func NewGrpcUserServer(grpc_endpoints GrpcUserServiceEndpoints) userpb.UserServiceServer {
	return &gRPCServer{
		createUser: grpc_gokit.NewServer(
			grpc_endpoints.CreateUser,
			decodeCreateUserRequest,
			encodeCreateUserResponse,
		),

		authenticate: grpc_gokit.NewServer(
			grpc_endpoints.Authenticate,
			decodeAuthenticateRequest,
			encodeAuthenticateResponse,
		),

		updateUser: grpc_gokit.NewServer(
			grpc_endpoints.UpdateUser,
			decodeUpdateUserRequest,
			encondeUpdatUserResponse,
		),
	}
}

func decodeCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	user_pb, ok := request.(*userpb.CreateUserRequest)

	if !ok {
		return nil, errors.New("No proto message 'CreateUserRequest' request")
	}

	req := CreateUserRequest{
		Email:     user_pb.GetEmail(),
		Password:  user_pb.GetPassword(),
		Age:       int(user_pb.GetAge()),
		ExtraInfo: user_pb.GetAdditionalInformation(),
	}

	return req, nil
}

func encodeCreateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(string)
	return &userpb.CreateUserResponse{Id: res}, nil
}

func decodeAuthenticateRequest(_ context.Context, request interface{}) (interface{}, error) {
	auth_pb, ok := request.(*userpb.AuthenticateRequest)

	if !ok {
		return nil, errors.New("No 'AuthenticateRequest' type")
	}

	req := AuthenticateRequest{
		Email:    auth_pb.GetEmail(),
		Password: auth_pb.GetPassword(),
	}

	return req, nil
}

func encodeAuthenticateResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(int)

	return &userpb.AuthenticateResponse{UserId: int32(res)}, nil
}

func decodeUpdateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	update_pb, ok := request.(*userpb.UpdateUserRequest)

	if !ok {
		return nil, errors.New("No 'UpdateUserRequest' type")
	}

	get_data := func(f string) interface{} {
		if f == "email" {
			return update_pb.GetEmail()
		} else if f == "password" {
			return update_pb.GetPassword()
		} else if f == "age" {
			return update_pb.GetAge()
		} else {
			return update_pb.GetAdditionalInformation()
		}
	}

	info := make(map[string]interface{})

	for _, field := range update_pb.GetFields() {
		info[field] = get_data(field)
	}

	req := UpdateUserRequest{
		Id:          int(update_pb.GetId()),
		Information: info,
	}

	return req, nil
}

func encondeUpdatUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(entities.User)

	return &userpb.UpdateUserResponse{
		Email:                 res.Email,
		Password:              res.Password,
		Age:                   uint32(res.Age),
		AdditionalInformation: res.ExtraInfo,
	}, nil
}

func (g *gRPCServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	_, res, err := g.createUser.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return res.(*userpb.CreateUserResponse), nil
}

func (g *gRPCServer) Authenticate(ctx context.Context, req *userpb.AuthenticateRequest) (*userpb.AuthenticateResponse, error) {
	_, res, err := g.authenticate.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return res.(*userpb.AuthenticateResponse), nil
}

func (g *gRPCServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	_, res, err := g.authenticate.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return res.(*userpb.UpdateUserResponse), nil
}
