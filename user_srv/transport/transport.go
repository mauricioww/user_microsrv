package transport

import (
	"context"
	"errors"

	grpc_gokit "github.com/go-kit/kit/transport/grpc"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
)

type gRPCServer struct {
	createUser   grpc_gokit.Handler
	authenticate grpc_gokit.Handler
	updateUser   grpc_gokit.Handler
	getUser      grpc_gokit.Handler
	deleteUser   grpc_gokit.Handler

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
			encondeUpdateUserResponse,
		),

		getUser: grpc_gokit.NewServer(
			grpc_endpoints.GetUser,
			decodeGetUserRequest,
			encodeGetUserResponse,
		),

		deleteUser: grpc_gokit.NewServer(
			grpc_endpoints.DeleteUser,
			decodeDeleteUserRequest,
			encondeDeleteUserResponse,
		),
	}
}

func decodeCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	create_pb, ok := request.(*userpb.CreateUserRequest)

	if !ok {
		return nil, errors.New("No proto message 'CreateUserRequest' request")
	}

	req := CreateUserRequest{
		Email:     create_pb.GetEmail(),
		Password:  create_pb.GetPassword(),
		Age:       int(create_pb.GetAge()),
		ExtraInfo: create_pb.GetAdditionalInformation(),
	}

	return req, nil
}

func encodeCreateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(CreateUserResponse)
	return &userpb.CreateUserResponse{Id: int32(res.Id)}, nil
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
	res := response.(AuthenticateResponse)

	return &userpb.AuthenticateResponse{UserId: int32(res.Id)}, nil
}

func decodeUpdateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	update_pb, ok := request.(*userpb.UpdateUserRequest)

	if !ok {
		return nil, errors.New("No 'UpdateUserRequest' type")
	}

	req := UpdateUserRequest{
		Id:        int(update_pb.GetId()),
		Email:     update_pb.GetEmail(),
		Password:  update_pb.GetPassword(),
		Age:       int(update_pb.GetAge()),
		ExtraInfo: update_pb.GetAdditionalInformation(),
	}

	return req, nil
}

func encondeUpdateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(UpdateUserResponse)

	return &userpb.UpdateUserResponse{
		Email:                 res.Email,
		Password:              res.Password,
		Age:                   uint32(res.Age),
		AdditionalInformation: res.ExtraInfo,
	}, nil
}

func decodeGetUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	get_pb, ok := request.(*userpb.GetUserRequest)

	if !ok {
		return nil, errors.New("No 'GetUserRequest' type")
	}

	req := GetUserRequest{
		UserId: int(get_pb.GetId()),
	}

	return req, nil
}

func encodeGetUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(GetUserResponse)

	return &userpb.GetUserResponse{
		Email:                 res.Email,
		Password:              res.Password,
		Age:                   uint32(res.Age),
		AdditionalInformation: res.ExtraInfo,
	}, nil
}

func decodeDeleteUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	delete_pb, ok := request.(*userpb.DeleteUserRequest)

	if !ok {
		return nil, errors.New("No 'DeleteUserRequest' type")
	}

	req := DeleteUserRequest{
		UserId: int(delete_pb.GetId()),
	}

	return req, nil
}

func encondeDeleteUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(DeleteUserResponse)

	return &userpb.DeleteUserResponse{Success: res.Success}, nil
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
	_, res, err := g.updateUser.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return res.(*userpb.UpdateUserResponse), nil
}

func (g *gRPCServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	_, res, err := g.getUser.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return res.(*userpb.GetUserResponse), nil
}

func (g *gRPCServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	_, res, err := g.deleteUser.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return res.(*userpb.DeleteUserResponse), nil
}
