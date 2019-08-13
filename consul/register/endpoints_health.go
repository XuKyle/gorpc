package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type HealthRequest struct {
}

type HealthResponse struct {
	Status bool `json:"status"`
}

func MakeHealthCheckEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		healthCheck := svc.HealthCheck()
		return HealthResponse{healthCheck}, nil
	}
}
