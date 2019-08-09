package main

import (
	"context"
	"encoding/json"
	kithttp "github.com/go-kit/kit/transport/http"
	"gorpc/official/ep"
	"gorpc/official/service/impl"
	"net/http"
)

func main() {
	stringService := impl.StringServiceImpl{}

	upperCaseHandler := kithttp.NewServer(ep.MakeUppercaseEndpoint(stringService), decodeUppercaseRequest, encodeResponse)
	countHandler := kithttp.NewServer(ep.MakeCountEndpoint(stringService), decodeCountRequest, encodeResponse)

	http.Handle("/uppercase", upperCaseHandler)
	http.Handle("/count", countHandler)

	http.ListenAndServe(":8080", nil)
}

// decode encode request
func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request ep.UppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request ep.CountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
