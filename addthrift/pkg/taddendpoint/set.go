package taddendpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"gorpc/addthrift/pkg/taddservice"
)

//New returns a Set that wraps the provided server, and wires in all of the
//expected endpoint middlewares via the various parameters.
func New(svc taddservice.Service, logger log.Logger, duration metrics.Histogram) Set {
	var sumEndpoint endpoint.Endpoint
	{
		sumEndpoint = MakeSumEndpoint(svc)
		sumEndpoint = LoggingMiddleware(log.With(logger, "method", "Sum"))(sumEndpoint)
		sumEndpoint = InstrumentingMiddleware(duration.With("method", "Sum"))(sumEndpoint)
	}

	var concatEndpoint endpoint.Endpoint
	{
		concatEndpoint = MakeConcatEndpoint(svc)
		concatEndpoint = LoggingMiddleware(log.With(logger, "emthod", "Concat"))(concatEndpoint)
		concatEndpoint = InstrumentingMiddleware(duration.With("method", "Concat"))(concatEndpoint)
	}

	return Set{
		sumEndpoint, concatEndpoint,
	}
}

// MakeSumEndpoint constructs a Sum endpoint wrapping the service.
func MakeSumEndpoint(service taddservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SumRequest)
		resp, err := service.Sum(ctx, req.A, req.B)
		return SumResponse{resp, err}, nil
	}
}

func MakeConcatEndpoint(service taddservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ConcatRequest)
		resp, err := service.Concat(ctx, req.A, req.B)
		return ConcatResponse{resp, err}, nil
	}
}

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.

type Set struct {
	SumEndpoint    endpoint.Endpoint
	ConcatEndpoint endpoint.Endpoint
}

// Sum implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) Sum(ctx context.Context, a, b int) (int, error) {
	resp, err := s.SumEndpoint(ctx, SumRequest{A: a, B: b})
	if err != nil {
		return 0, err
	}
	response := resp.(SumResponse)
	return response.V, response.Err
}
func (s Set) Concat(ctx context.Context, a, b string) (string, error) {
	resp, err := s.ConcatEndpoint(ctx, ConcatRequest{a, b})
	if err != nil {
		return "", err
	}
	response := resp.(ConcatResponse)
	return response.V, response.Err
}

// request and response
// SumRequest collects the request parameters for the Sum method.
type SumRequest struct {
	A, B int
}

// SumResponse collects the response values for the Sum method.
type SumResponse struct {
	V   int   `json:"v"`
	Err error `json:"-"` // should be intercepted by Failed/errorEncoder
}

// Failed implements endpoint.Failer.
func (r SumResponse) Failed() error { return r.Err }

// ConcatRequest collects the request parameters for the Concat method.
type ConcatRequest struct {
	A, B string
}

// ConcatResponse collects the response values for the Concat method.
type ConcatResponse struct {
	V   string `json:"v"`
	Err error  `json:"-"`
}

// Failed implements endpoint.Failer.
func (r ConcatResponse) Failed() error { return r.Err }
