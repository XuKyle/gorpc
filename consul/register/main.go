package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
)

func main() {

	var (
		//consulHost  = flag.String("consul.host", "", "consul ip address")
		//consulPort  = flag.String("consul.port", "", "consul port")
		//serviceHost = flag.String("service.host", "", "service ip address")
		servicePort = flag.String("service.port", "9000", "service port")
	)

	flag.Parse()

	var svc Service
	svc = ArithmeticService{}

	ctx := context.Background()
	errChan := make(chan error)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	endpoint := MakeArithmeticEndpoint(svc)
	healthEndpoint := MakeHealthCheckEndpoint(svc)

	serverEndpoints := ArithmeticEndpoints{
		ArithmeticEndpoint:  endpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	//创建http.Handler
	r := MakeHttpHandler(ctx, serverEndpoints, logger)

	go func() {
		fmt.Println("Http Server start at port:" + *servicePort)
		handler := r
		errChan <- http.ListenAndServe(":"+*servicePort, handler)
	}()

	error := <-errChan
	fmt.Println(error)
}
