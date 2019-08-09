package ep

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"gorpc/official/service"
)

// request response
type UppercaseRequest struct {
	S string `json:"s"`
}

type UppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

type CountRequest struct {
	S string `json:"s"`
}

type CountResponse struct {
	V int `json:"v"`
}

// 添加endpoint
func MakeUppercaseEndpoint(service service.StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		uppercaseRequest := request.(UppercaseRequest)
		upperCase, err := service.UpperCase(uppercaseRequest.S)
		if err != nil {
			return UppercaseResponse{upperCase, err.Error()}, err
		}
		return UppercaseResponse{upperCase, ""}, nil
	}
}

func MakeCountEndpoint(stringService service.StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		countRequest := request.(CountRequest)
		count := stringService.Count(countRequest.S)
		return CountResponse{count}, nil
	}
}
