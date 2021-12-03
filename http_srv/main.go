package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/http_srv/repository"
	"github.com/mauricioww/user_microsrv/http_srv/service"
	"github.com/mauricioww/user_microsrv/http_srv/transport"
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

	ctx := context.Background()
	var http_srv service.HttpService
	{
		repository := repository.NewHttpRepository(logger)
		http_srv = service.NewHttpService(repository, logger)
	}

	err := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		err <- fmt.Errorf("%s", <-c)
	}()

	http_endpoints := transport.MakeHttpEndpoints(http_srv)

	go func() {
		fmt.Println("Listengin on port: 8080")
		http_handler := transport.NewHTTPServer(ctx, http_endpoints)
		err <- http.ListenAndServe(":8080", http_handler)
	}()

	level.Error(logger).Log("exit: ", <-err)
}
