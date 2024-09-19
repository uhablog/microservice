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
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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
		Price:  3000,
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

func (s *server) ListBooks(ctx context.Context, in *emptypb.Empty) (*pb.ListBooksResponse, error) {
	protoBooks := make([]*pb.Book, 0)

	for _, book := range books {
		protoBook := &pb.Book{
			Id:     int32(book.Id),
			Title:  book.Title,
			Author: book.Author,
			Price:  int32(book.Price),
		}
		protoBooks = append(protoBooks, protoBook)
	}

	return &pb.ListBooksResponse{Books: protoBooks}, nil
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
