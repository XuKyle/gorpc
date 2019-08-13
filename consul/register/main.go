package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	var (
		//consulHost  = flag.String("consul.host", "", "consul ip address")
		//consulPort  = flag.String("consul.port", "", "consul port")
		//serviceHost = flag.String("service.host", "", "service ip address")
		servicePort = flag.String("service.port", "9000", "service port")
	)

	flag.Parse()

	// *************logger***************
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// *************metric***************
	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "raysonxin",
		Subsystem: "arithmetic_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "raysonxin",
		Subsystem: "arithemetic_service",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	var svc Service
	svc = ArithmeticService{}

	// 添加middleware
	svc = LoggingMiddleware(logger)(svc)
	svc = Metrics(requestCount, requestLatency)(svc)

	ctx := context.Background()
	errChan := make(chan error)

	// ************************************ endpoint ************************************
	endpoint := MakeArithmeticEndpoint(svc)
	// add ratelimit,refill every second,set capacity 3
	//ratebucket := ratelimit.NewBucket(time.Second*1, 3)
	//endpoint = NewTokenBucketLimitterWithJuju(ratebucket)(endpoint)

	//add ratelimit,refill every second,set capacity 3
	ratebucket := rate.NewLimiter(rate.Every(time.Second*1), 100)
	endpoint = NewTokenBucketLimitterWithBuildIn(ratebucket)(endpoint)

	healthEndpoint := MakeHealthCheckEndpoint(svc)

	serverEndpoints := ArithmeticEndpoints{
		ArithmeticEndpoint:  endpoint,
		HealthCheckEndpoint: healthEndpoint,
	}
	// ************************************ endpoint ************************************

	//创建http.Handler
	r := MakeHttpHandler(ctx, serverEndpoints, logger)

	go func() {
		fmt.Println("Http Server start at port:" + *servicePort)
		handler := r
		errChan <- http.ListenAndServe(":"+*servicePort, handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	fmt.Println(error)
}
