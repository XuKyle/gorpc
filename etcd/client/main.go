package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	"google.golang.org/grpc"
	"gorpc/book"
	"io"
	"time"
)

func main() {
	var (
		etcdServer = "127.0.0.1:2379"
		prefix     = "/services/book/"
		ctx        = context.Background()
	)

	clientOptions := etcdv3.ClientOptions{DialTimeout: time.Second * 3, DialKeepAlive: time.Second * 3}
	// 注册中心
	client, err := etcdv3.NewClient(ctx, []string{etcdServer}, clientOptions)
	if err != nil {
		panic(err)
	}

	logger := log.NewNopLogger()

	//创建实例管理器, 此管理器会Watch监听etc中prefix的目录变化更新缓存的服务实例数据
	instancer, err := etcdv3.NewInstancer(client, prefix, logger)
	if err != nil {
		panic(err)
	}
	//创建端点管理器， 此管理器根据Factory和监听的到实例创建endPoint并订阅instancer的变化动态更新Factory创建的endPoint
	endpointer := sd.NewEndpointer(instancer, reqFactory, logger)

	//创建负载均衡器
	balancer := lb.NewRoundRobin(endpointer)

	/**
	  我们可以通过负载均衡器直接获取请求的endPoint，发起请求
	  reqEndPoint,_ := balancer.Endpoint()
	*/

	/**
	  也可以通过retry定义尝试次数进行请求
	*/
	reqEndPoint := lb.Retry(3, 3*time.Second, balancer)

	req := struct{}{}

	// 单个请求版本
	//if _, err := reqEndPoint(ctx, req); err != nil {
	//	panic(nil)
	//}

	for i := 0; i < 10; i++ {
		if _, err := reqEndPoint(ctx, req); err != nil {
			fmt.Println("error ")
		}
	}
}

func reqFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		fmt.Println("请求服务:", instanceAddr, "当前时间：", time.Now().Format("2006-01-02 15:04:05"))
		conn, err := grpc.Dial(instanceAddr, grpc.WithInsecure())
		if err != nil {
			fmt.Println(err)
		}
		defer conn.Close()

		bookServiceClient := book.NewBookServiceClient(conn)
		bookInfo, err := bookServiceClient.GetBookInfo(context.Background(), &book.BookInfoParams{BookId: 1})
		fmt.Println("bookInfo:")
		fmt.Println("book：1 =>", bookInfo)

		bookList, err := bookServiceClient.GetBookList(context.Background(), &book.BookListParams{Page: 1, Limit: 2})
		fmt.Println("bookList:")
		for _, list := range bookList.BookList {
			fmt.Println(list)
		}

		return nil, nil
	}, nil, nil
}
