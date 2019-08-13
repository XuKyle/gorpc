package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"gorpc/addthrift/pkg/taddservice"
	"gorpc/addthrift/pkg/taddtransport"
	"gorpc/addthrift/thrift/gen-go/addsvc"
	"os"
	"strconv"
	"text/tabwriter"
)

func main() {
	fs := flag.NewFlagSet("addcli", flag.ExitOnError)

	var thriftAddr = fs.String("thrift-addr", ":8083", "Thrift address of addsvc")
	var method = fs.String("method", "sum", "sum, concat")

	fs.Usage = usageForCli(fs, os.Args[0]+" [flags] <a> <b>")
	fs.Parse(os.Args[1:])
	if len(fs.Args()) != 2 {
		fs.Usage()
		os.Exit(1)
	}

	var (
		svc taddservice.Service
		err error
	)

	var protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTTransportFactory()
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)

	transportSocket, err := thrift.NewTSocket(*thriftAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "socket error:%v", err)
	}

	transport, err := transportFactory.GetTransport(transportSocket)
	if err := transport.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer transport.Close()
	client := addsvc.NewAddServiceClientFactory(transport, protocolFactory)
	svc = taddtransport.NewThriftClient(client)

	switch *method {
	case "sum":
		a, _ := strconv.ParseInt(fs.Args()[0], 10, 64)
		b, _ := strconv.ParseInt(fs.Args()[1], 10, 64)
		v, err := svc.Sum(context.Background(), int(a), int(b))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%d + %d = %d\n", a, b, v)

	case "concat":
		a := fs.Args()[0]
		b := fs.Args()[1]
		v, err := svc.Concat(context.Background(), a, b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%q + %q = %q\n", a, b, v)

	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", *method)
		os.Exit(1)
	}

}

func usageForCli(fs *flag.FlagSet, short string) func() {
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
