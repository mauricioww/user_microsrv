package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/user_details_srv/detailspb"
	"github.com/mauricioww/user_microsrv/user_details_srv/repository"
	"github.com/mauricioww/user_microsrv/user_details_srv/service"
	"github.com/mauricioww/user_microsrv/user_details_srv/transport"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

const (
	mongo_uri = "mongodb://127.0.0.1:27017"
	database  = "gokitexa"
	user_name = "mauriciow"
	password  = "password"
)

func main() {

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(
			logger,
			"service",
			"GRPC_USER_DETAILS",
			"time",
			log.DefaultTimestampUTC,
			"caller",
			log.DefaultCaller,
		)
	}
	level.Info(logger).Log("mesg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	var db *mongo.Database

	{
		credentials := options.Credential{
			Username: user_name,
			Password: password,
		}
		client_opts := options.Client().ApplyURI(mongo_uri).SetAuth(credentials)
		client, err := mongo.Connect(context.Background(), client_opts)

		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}

		db = client.Database(database)
	}

	var grpc_user_details_srv service.GrpcUserDetailsService
	{
		mongo_repository := repository.NewUserDetailsRepository(db, logger)
		grpc_user_details_srv = service.NewGrpcUserDetailsService(mongo_repository, logger)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	grpc_endpoints := transport.MakeGrpcUserDetailsServiceEndpoints(grpc_user_details_srv)
	grpc_server := transport.NewGrpcUserDetailsServer(grpc_endpoints)
	grpc_listener, err := net.Listen("tcp", ":50052")

	if err != nil {
		logger.Log("Error listening: ", err)
		os.Exit(-1)
	}

	go func() {
		server := grpc.NewServer()
		detailspb.RegisterUserDetailsServiceServer(server, grpc_server)
		if err := server.Serve(grpc_listener); err != nil {
			logger.Log("Error serving", err)
		}
		level.Info(logger).Log("info", "grpc server started")
	}()

	level.Error(logger).Log("exit: ", <-errs)
}
