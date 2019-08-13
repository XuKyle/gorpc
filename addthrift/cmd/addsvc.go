package main

import (
	"flag"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorpc/addthrift/pkg/taddendpoint"
	"gorpc/addthrift/pkg/taddservice"
	"gorpc/addthrift/pkg/taddtransport"
	"gorpc/addthrift/thrift/gen-go/addsvc"
	"gorpc/addthrift/utils"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"
)

func main() {
	fs := flag.NewFlagSet("addservice", flag.ExitOnError)
	thriftAddr := fs.String("thrift-addr", ":8083", "thrift listen address")

	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	_ = fs.Parse(os.Args[1:])

	// Create a single logger, which we'll use and give to other components.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Create the (sparse) metrics we'll use in the service. They, too, are
	// dependencies that we pass to components that use them.
	var ints, chars metrics.Counter
	{
		// Business-level metrics.
		ints = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "example",
			Subsystem: "addsvc",
			Name:      "integers_summed",
			Help:      "Total count of integers summed via the Sum method.",
		}, []string{})
		chars = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "example",
			Subsystem: "addsvc",
			Name:      "characters_concatenated",
			Help:      "Total count of characters concatenated via the Concat method.",
		}, []string{})
	}
	var duration metrics.Histogram
	{
		// Endpoint-level metrics.
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "example",
			Subsystem: "addsvc",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
	}
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	// server
	service := taddservice.New(logger, ints, chars)
	endpoint := taddendpoint.New(service, logger, duration)
	thriftServer := taddtransport.NewThriftServer(endpoint)

	// The Thrift socket mounts the Go kit Thrift server we created earlier.
	// There's a lot of boilerplate involved here, related to configuring
	// the protocol and transport; blame Thrift.
	thriftSocket, err := thrift.NewTServerSocket(*thriftAddr)
	if err != nil {
		logger.Log("transport", "Thrift", "during", "Listen", "err", err)
		os.Exit(1)
	}

	logger.Log("transport", "Thrift", "addr", *thriftAddr)
	var protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	var transportFactory = thrift.NewTTransportFactory()
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)

	err = thrift.NewTSimpleServer4(
		addsvc.NewAddServiceProcessor(thriftServer),
		thriftSocket,
		transportFactory,
		protocolFactory,
	).Serve()

	// 注册
	registar := Register("192.168.70.3", "8500", utils.GetLocalIp(), "8083", logger)
	registar.Register()

	if err != nil {
		fmt.Println("server error")
	}
}

func Register(consulHost, consulPort, svcHost, svcPort string, logger log.Logger) (registar sd.Registrar) {
	fmt.Println("register***********")
	// 创建Consul客户端连接
	var client consul.Client
	{
		consulCfg := api.DefaultConfig()
		consulCfg.Address = consulHost + ":" + consulPort
		consulClient, err := api.NewClient(consulCfg)
		if err != nil {
			logger.Log("create consul client error:", err)
			os.Exit(1)
		}

		client = consul.NewClient(consulClient)
	}

	// 设置Consul对服务健康检查的参数
	check := api.AgentServiceCheck{
		HTTP:     "http://" + svcHost + ":" + svcPort + "/health",
		Interval: "10s",
		Timeout:  "1s",
		Notes:    "Consul check service health status.",
	}

	port, _ := strconv.Atoi(svcPort)

	//设置微服务想Consul的注册信息
	reg := api.AgentServiceRegistration{
		ID:      "add" + svcHost,
		Name:    "addsvc",
		Address: svcHost,
		Port:    port,
		Tags:    []string{"recommend", "add"},
		Check:   &check,
	}

	// 执行注册
	registar = consul.NewRegistrar(client, &reg, logger)
	return
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}
