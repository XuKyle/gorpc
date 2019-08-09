package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"gorpc/book"
)

func main() {
	serviceAddr := "127.0.0.1:9000"
	conn, err := grpc.Dial(serviceAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("server error!")
	}

	defer conn.Close()

	client := book.NewBookServiceClient(conn)
	bookInfo, _ := client.GetBookInfo(context.Background(), &book.BookInfoParams{BookId: 1})
	fmt.Println("*** 书籍 ***")
	fmt.Println(bookInfo)

	bookList, err := client.GetBookList(context.Background(), &book.BookListParams{})
	fmt.Println("*** 书籍列表 ***")
	for k, v := range bookList.BookList {
		fmt.Printf("%d : %s\n", k, v)
	}

}
