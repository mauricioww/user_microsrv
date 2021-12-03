package transport

import (
	"context"
	"errors"
	"log"

	grpc_gokit "github.com/go-kit/kit/transport/grpc"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
)

type gRPCServer struct {
	createUser grpc_gokit.Handler
	userpb.UnimplementedUserServiceServer
}

func NewGrpcUserServer(grpc_endpoints GrpcUserServiceEndpoints, logger log.Logger) userpb.UserServiceServer {
	return &gRPCServer{
		createUser: grpc_gokit.NewServer(
			grpc_endpoints.CreateUser,
			decodeCreateUserRequest,
			encodeCreateUserResponse,
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
		Password:  user_pb.GetEmail(),
		Age:       int(user_pb.GetAge()),
		ExtraInfo: user_pb.GetAdditionalInformation(),
	}

	return req, nil
}

func encodeCreateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(CreateUserResponse)
	return &userpb.CreateUserResponse{Id: res.Id}, nil
}

func (g *gRPCServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	_, res, err := g.createUser.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return res.(*userpb.CreateUserResponse), nil
}
