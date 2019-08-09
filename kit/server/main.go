package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kit_grpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"gorpc/book"
	"net"
)

type BookServer struct {
	bookListHandler kit_grpc.Handler
	bookInfoHandler kit_grpc.Handler
}

//通过grpc调用GetBookInfo时,GetBookInfo只做数据透传, 调用BookServer中对应Handler.ServeGRPC转交给go-kit处理
func (server *BookServer) GetBookInfo(ctx context.Context, in *book.BookInfoParams) (*book.BookInfo, error) {
	_, resp, e := server.bookInfoHandler.ServeGRPC(ctx, in)
	if e != nil {
		fmt.Println("server happen!")
	}
	return resp.(*book.BookInfo), nil
}

//通过grpc调用GetBookList时,GetBookList只做数据透传, 调用BookServer中对应Handler.ServeGRPC转交给go-kit处理
func (server *BookServer) GetBookList(ctx context.Context, in *book.BookListParams) (*book.BookList, error) {
	_, rsp, err := server.bookListHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return rsp.(*book.BookList), err
}

//创建bookList的EndPoint 
func makeGetBookListEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//请求列表时返回 书籍列表
		bl := new(book.BookList)
		bl.BookList = append(bl.BookList, &book.BookInfo{BookId: 1, BookName: "go in action"})
		bl.BookList = append(bl.BookList, &book.BookInfo{BookId: 2, BookName: "springboot in action"})
		bl.BookList = append(bl.BookList, &book.BookInfo{BookId: 3, BookName: "kotlin in action"})
		return bl, nil
	}
}

//创建bookInfo的EndPoint
func makeGetBookInfoEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//请求详情时返回 书籍信息
		req := request.(*book.BookInfoParams)
		b := new(book.BookInfo)
		b.BookId = req.BookId
		b.BookName = "go go go!"
		return b, nil
	}
}

func decodeRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func encodeRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func main() {
	// 创建bookInfo的Handler
	bookInfoHandler := kit_grpc.NewServer(makeGetBookInfoEndpoint(), decodeRequest, encodeRequest)
	// 创建bookList的Handler
	bookListHandler := kit_grpc.NewServer(makeGetBookListEndpoint(), decodeRequest, encodeRequest)

	//包装BookServer
	bookServer := new(BookServer)

	//bookServer 增加 go-kit流程的 bookInfo处理逻辑   增加 go-kit流程的 bookList处理逻辑
	bookServer.bookInfoHandler = bookInfoHandler
	bookServer.bookListHandler = bookListHandler

	//启动grpc服务
	serverAddr := ":9000"
	listener, _ := net.Listen("tcp", serverAddr)
	gserver := grpc.NewServer()
	book.RegisterBookServiceServer(gserver, bookServer)
	err := gserver.Serve(listener)
	if err != nil {
		fmt.Println("server error")
	}
}
