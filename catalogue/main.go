package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "gihyo/catalogue/proto/book"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Book struct {
	Id     int
	Title  string
	Author string
	Price  int
}

var (
	book1 = Book{
		Id:     1,
		Title:  "The Awakening",
		Author: "Kate Chopin",
		Price:  1000,
	}
	book2 = Book{
		Id:     2,
		Title:  "City of Glass",
		Author: "Paul Auster",
		Price:  2000,
	}
	books = []Book{book1, book2}
)

func getBook(i int32) Book {
	return books[i-1]
}

type server struct {
	pb.UnimplementedCatalogueServer
}

func (s *server) GetBook(ctx context.Context, in *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	book := getBook(in.Id)

	protoBook := &pb.Book{
		Id:     int32(book.Id),
		Title:  book.Title,
		Author: book.Author,
		Price:  int32(book.Price),
	}

	return &pb.GetBookResponse{Book: protoBook}, nil
}

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCatalogueServer(s, &server{})
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
