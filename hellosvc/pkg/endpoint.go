package pkg

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	HelloEndpoint endpoint.Endpoint
}

func NewEndpoint(svc Service) Endpoints {
	var he endpoint.Endpoint
	{
		he = MakeHelloEndpoint(svc)
	}

	return Endpoints{
		HelloEndpoint: he,
	}
}

func MakeHelloEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(HelloRequest)
		v, err := s.Hello(ctx, req.N)
		return HelloResponse{H: v, Err: err}, nil
	}
}

type HelloRequest struct {
	N string
}

type HelloResponse struct {
	H string
	Err error
}