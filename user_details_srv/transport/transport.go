package transport

import (
	"context"
	"errors"

	grpc_gokit "github.com/go-kit/kit/transport/grpc"
	"github.com/mauricioww/user_microsrv/user_details_srv/detailspb"
)

type gRPCServer struct {
	setUserDetails grpc_gokit.Handler

	detailspb.UnimplementedUserDetailsServiceServer
}

func NewGrpcUserDetailsServer(endpoints GrpcUserDetailsServiceEndpoints) detailspb.UserDetailsServiceServer {
	return &gRPCServer{
		setUserDetails: grpc_gokit.NewServer(
			endpoints.SetUserDetails,
			decodeSetUserDetailsRequest,
			encodeSetUserDetailsResponse,
		),
	}
}

func decodeSetUserDetailsRequest(_ context.Context, request interface{}) (interface{}, error) {
	set_details_pb, ok := request.(*detailspb.SetUserDetailsRequest)

	if !ok {
		return nil, errors.New("No proto message 'SetUserDetailsRequest' request")
	}

	req := SetUserDetailsRequest{
		UserId:       int(set_details_pb.GetUserId()),
		Country:      set_details_pb.GetCountry(),
		City:         set_details_pb.GetCity(),
		MobileNumber: set_details_pb.GetMobileNumber(),
		Married:      set_details_pb.GetMarried(),
		Height:       set_details_pb.GetHeight(),
		Weigth:       set_details_pb.Weigth,
	}

	return req, nil
}

func encodeSetUserDetailsResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(SetUserDetailsResponse)
	return &detailspb.SetUserDetailsResponse{Success: res.Success}, nil
}

func (g *gRPCServer) SetUserDetails(ctx context.Context, req *detailspb.SetUserDetailsRequest) (*detailspb.SetUserDetailsResponse, error) {
	_, res, err := g.setUserDetails.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return res.(*detailspb.SetUserDetailsResponse), nil
}
