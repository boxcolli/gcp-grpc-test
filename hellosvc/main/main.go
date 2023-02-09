package main

import (
	"flag"
	"fmt"
	"hellosvc/pkg"
	"net"
	"os"
	"os/signal"
	"syscall"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"google.golang.org/grpc"

	"hellosvc/pb"
)

func main() {
	var (
		grpcAddr = flag.String("grpc-addr", ":8082", "gRPC listen address")		
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var (
		service = pkg.NewService(logger)
		endpoints = pkg.NewEndpoint(service)
		grpcServer = pkg.NewGRPCServer(endpoints, logger)
	)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	{
		grpcListener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
		pb.RegisterHelloServer(baseServer, grpcServer)
		errs <- baseServer.Serve(grpcListener)
	}	

	logger.Log("exit", <-errs)
}
