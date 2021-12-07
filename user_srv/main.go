package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mauricioww/user_microsrv/user_srv/repository"
	"github.com/mauricioww/user_microsrv/user_srv/service"
	"github.com/mauricioww/user_microsrv/user_srv/transport"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
	"google.golang.org/grpc"
)

const (
	mysql_database = "mauricio:password@tcp(127.0.0.1:3306)/gokitexa"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(
			logger,
			"service",
			"HTTP_SRV",
			"time",
			log.DefaultTimestampUTC,
			"caller",
			log.DefaultCaller,
		)
	}

	level.Info(logger).Log("mesg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	var db *sql.DB
	{
		var err error
		db, err = sql.Open("mysql", mysql_database)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}

	var grpc_user_srv service.GrpcUserService
	{
		mysql_repository := repository.NewUserSrvRepository(db, logger)
		grpc_user_srv = service.NewGrpcUserService(mysql_repository, logger)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	grpc_endpoints := transport.MakeGrpcUserServiceEndpoints(grpc_user_srv)
	grpc_server := transport.NewGrpcUserServer(grpc_endpoints)
	grpc_listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		logger.Log("Error listening: ", err)
		os.Exit(-1)
	}

	go func() {
		server := grpc.NewServer()
		userpb.RegisterUserServiceServer(server, grpc_server)
		if err := server.Serve(grpc_listener); err != nil {
			logger.Log("Error serving", err)
		}
		level.Info(logger).Log("info", "grpc server started")
	}()

	level.Error(logger).Log("exit: ", <-errs)
}