package taddtransport

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"gorpc/addthrift/pkg/taddendpoint"
	"gorpc/addthrift/pkg/taddservice"
	"gorpc/addthrift/thrift/gen-go/addsvc"
)

// NewThriftServer makes a set of endpoints available as a Thrift service.
func NewThriftServer(endpoints taddendpoint.Set) addsvc.AddService {
	return &thriftServer{
		endpoints: endpoints,
	}
}

type thriftServer struct {
	ctx       context.Context
	endpoints taddendpoint.Set
}

func (s *thriftServer) Sum(ctx context.Context, a int64, b int64) (*addsvc.SumReply, error) {
	request := taddendpoint.SumRequest{int(a), int(b)}
	response, err := s.endpoints.SumEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	resp := response.(taddendpoint.SumResponse)
	return &addsvc.SumReply{Value: int64(resp.V), Err: err2str(resp.Err)}, nil
}

func (s *thriftServer) Concat(ctx context.Context, a string, b string) (r *addsvc.ConcatReply, err error) {
	request := taddendpoint.ConcatRequest{a, b}
	response, err := s.endpoints.ConcatEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	resp := response.(taddendpoint.ConcatResponse)
	return &addsvc.ConcatReply{Value: resp.V, Err: err2str(resp.Err)}, nil
}

// client *********
func NewThriftClient(client *addsvc.AddServiceClient) taddservice.Service {
	var sumEndpoint endpoint.Endpoint = MakeThriftSumEndpoint(client)
	var concatEndpoint endpoint.Endpoint = MakeThriftConcatEndpoint(client)

	return taddendpoint.Set{
		SumEndpoint:    sumEndpoint,
		ConcatEndpoint: concatEndpoint,
	}
}

func MakeThriftSumEndpoint(client *addsvc.AddServiceClient) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(taddendpoint.SumRequest)
		reply, err := client.Sum(ctx, int64(req.A), int64(req.B))
		return taddendpoint.SumResponse{V: int(reply.Value), Err: err}, nil
	}
}

func MakeThriftConcatEndpoint(client *addsvc.AddServiceClient) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(taddendpoint.ConcatRequest)
		reply, err := client.Concat(ctx, req.A, req.B)
		return taddendpoint.ConcatResponse{reply.Value, err}, nil
	}
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
