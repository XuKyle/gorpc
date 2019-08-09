package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"
	kit_grpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"gorpc/book"
	"net"
	"time"
)

type BookServer struct {
	bookListHandler kit_grpc.Handler
	bookInfoHandler kit_grpc.Handler
}

func (server *BookServer) GetBookInfo(ctx context.Context, in *book.BookInfoParams) (*book.BookInfo, error) {
	_, resp, err := server.bookInfoHandler.ServeGRPC(ctx, in)
	if err != nil {
		fmt.Println("getBookInfo error")
	}
	return resp.(*book.BookInfo), nil
}

func (server *BookServer) GetBookList(ctx context.Context, in *book.BookListParams) (*book.BookList, error) {
	_, resp, err := server.bookListHandler.ServeGRPC(ctx, in)
	if err != nil {
		fmt.Println("getBookList error")
	}
	return resp.(*book.BookList), nil
}

func makeGetBookInfoEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		bookInfo := new(book.BookInfo)
		bookInfo.BookId = 1
		bookInfo.BookName = "go in action!"
		return bookInfo, nil
	}
}

//创建bookList的EndPoint
func makeGetBookListEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//请求列表时返回 书籍列表
		bl := new(book.BookList)
		bl.BookList = append(bl.BookList, &book.BookInfo{BookId: 1, BookName: "springboot"})
		bl.BookList = append(bl.BookList, &book.BookInfo{BookId: 2, BookName: "go"})
		bl.BookList = append(bl.BookList, &book.BookInfo{BookId: 3, BookName: "kotlin"})
		bl.BookList = append(bl.BookList, &book.BookInfo{BookId: 4, BookName: "java"})
		return bl, nil
	}
}

func decodeRequest(_ context.Context, req interface{}) (interface{}, error) {
	return req, nil
}

func encodeResponse(_ context.Context, rsp interface{}) (interface{}, error) {
	return rsp, nil
}

func main() {
	var (
		//etcd服务地址
		etcdServer = "127.0.0.1:2379"
		//服务的信息目录
		prefix = "/services/book/"

		// 多个实例切换地址
		//当前启动服务实例的地址
		instance = "127.0.0.1:50051"
		//服务监听地址
		serviceAddress = ":50051"

		//服务实例注册的路径
		key = prefix + instance
		//服务实例注册的val
		value = instance
		ctx   = context.Background()
	)

	//etcd的连接参数
	clientOptions := etcdv3.ClientOptions{DialTimeout: time.Second * 3, DialKeepAlive: time.Second * 3}
	// 创建etcd链接
	client, err := etcdv3.NewClient(ctx, []string{etcdServer}, clientOptions)
	if err != nil {
		panic(err)
	}

	// 服务
	service := etcdv3.Service{Key: key, Value: value}
	// 注册
	registrar := etcdv3.NewRegistrar(client, service, log.NewNopLogger())

	// 启动注册
	registrar.Register()

	// 服务器
	bookInfoHandler := kit_grpc.NewServer(makeGetBookInfoEndpoint(), decodeRequest, encodeResponse)
	bookListHandler := kit_grpc.NewServer(makeGetBookListEndpoint(), decodeRequest, encodeResponse)

	bookServer := new(BookServer)
	bookServer.bookInfoHandler = bookInfoHandler
	bookServer.bookListHandler = bookListHandler

	listener, _ := net.Listen("tcp", serviceAddress)
	gserver := grpc.NewServer(grpc.UnaryInterceptor(kit_grpc.Interceptor))
	book.RegisterBookServiceServer(gserver, bookServer)

	err = gserver.Serve(listener)
	if err != nil {
		fmt.Println("error happen!")
	}

}
