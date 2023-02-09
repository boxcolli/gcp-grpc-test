package pkg

import (
	"context"

	"github.com/go-kit/log"
)

type Service interface {
	Hello(ctx context.Context, name string) (string, error)
}

func NewService(logger log.Logger) Service {
	var svc Service
	{
		svc = NewBasicService()
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

type basicService struct {

}

func NewBasicService() Service {
	return basicService{}
}

func (s basicService) Hello(ctx context.Context, name string) (string, error) {
	return ("Hello, " + name + "!"), nil
}

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) Hello(ctx context.Context, name string) (v string, err error) {
	defer func() {
		mw.logger.Log("method", "Hello", "name", name, "v", v, "err", err)
	}()
	return mw.next.Hello(ctx, name)
}