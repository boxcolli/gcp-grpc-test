package pkg

import (
	"context"

	"github.com/go-kit/kit/transport"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	"hellosvc/pb"
)

type grpcServer struct {
	hello grpctransport.Handler
	pb.UnimplementedHelloServer
}

func NewGRPCServer(endpoints Endpoints, logger log.Logger) pb.HelloServer{
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		hello: grpctransport.NewServer(
			endpoints.HelloEndpoint,
			decodeGRPCHelloRequest,
			encodeGRPCHelloResponse,
			append(options)...,
		),
	}
}

func (s *grpcServer) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	_, rep, err := s.hello.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.HelloReply), nil
}

// grpc request -> user-domain request
func decodeGRPCHelloRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.HelloRequest)
	return HelloRequest{N: req.Name}, nil
}

func encodeGRPCHelloResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	res := grpcReply.(HelloResponse)
	return &pb.HelloReply{Hello: res.H, Error: err2str(res.Err)}, nil
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
