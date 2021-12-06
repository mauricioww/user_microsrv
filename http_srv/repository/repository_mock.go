package repository

import (
	"context"
	"log"
	"net"

	"github.com/mauricioww/user_microsrv/user_srv/userpb"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type GrpcMock struct {
	mock.Mock
	userpb.UnimplementedUserServiceServer
}

func Dialer(gm *GrpcMock) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	userpb.RegisterUserServiceServer(server, gm.UnimplementedUserServiceServer)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func (m *GrpcMock) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	args := m.Called(ctx, req)

	return args.Get(0).(*userpb.CreateUserResponse), args.Error(1)
}
