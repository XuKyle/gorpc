package taddservice

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

var (
	// ErrTwoZeroes is an arbitrary business rule for the Add method.
	ErrTwoZeroes = errors.New("can't sum two zeroes")

	// ErrIntOverflow protects the Add method. We've decided that this error
	// indicates a misbehaving service and should count against e.g. circuit
	// breakers. So, we return it directly in endpoints, to illustrate the
	// difference. In a real service, this probably wouldn't be the case.
	ErrIntOverflow = errors.New("integer overflow")

	// ErrMaxSizeExceeded protects the Concat method.
	ErrMaxSizeExceeded = errors.New("result exceeds maximum size")
)

const (
	intMax = 1<<31 - 1
	intMin = -(intMax + 1)
	maxLen = 10
)

// Service describes a service that adds things together.
type Service interface {
	Sum(ctx context.Context, a, b int) (int, error)
	Concat(ctx context.Context, a, b string) (string, error)
}

// 定义service的实现
type basicService struct{}

func (basicService) Sum(ctx context.Context, a, b int) (int, error) {
	if a == 0 && b == 0 {
		return 0, ErrTwoZeroes
	}
	if b > 0 && a > (intMax-b) || (b < 0 && a < (intMin-b)) {
		return 0, ErrIntOverflow
	}
	return a + b, nil
}

func (basicService) Concat(ctx context.Context, a, b string) (string, error) {
	if len(a)+len(b) > maxLen {
		return "", ErrMaxSizeExceeded
	}
	return a + b, nil
}

// NewBasicService returns a naïve, stateless implementation of Service.
func NewBasicService() Service {
	return basicService{}
}

// New returns a basic Service with all of the expected middlewares wired in.
func New(logger log.Logger, ints, chars metrics.Counter) Service {
	var svc Service
	{
		svc = NewBasicService()
		svc = LoggingMiddleware(logger)(svc)
		svc = InstrumentingMiddleware(ints, chars)(svc)
	}
	return svc
}
