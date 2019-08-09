package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gorpc/book"
	"net"
)

type BookServer struct{}

func (*BookServer) GetBookInfo(ctx context.Context, in *book.BookInfoParams) (*book.BookInfo, error) {
	bookInfo := new(book.BookInfo)
	bookInfo.BookId = in.BookId
	bookInfo.BookName = "go in action"
	return bookInfo, nil
}

func (*BookServer) GetBookList(context.Context, *book.BookListParams) (*book.BookList, error) {
	bookList := new(book.BookList)
	bookList.BookList = append(bookList.BookList, &book.BookInfo{BookId: 1, BookName: "go in action!"})
	bookList.BookList = append(bookList.BookList, &book.BookInfo{BookId: 2, BookName: "spring in action"})
	return bookList, nil
}

func main() {
	serviceAddress := ":9000"
	bookServer := new(BookServer)
	//创建tcp监听
	ls, _ := net.Listen("tcp", serviceAddress)
	//创建grpc服务
	gs := grpc.NewServer()
	//注册bookServer
	book.RegisterBookServiceServer(gs, bookServer)
	//启动服务
	gs.Serve(ls)
}