package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/gdochadipa/oauth2-go-project/pkg/pb_test"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb_test.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb_test.HelloRequest) (*pb_test.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb_test.HelloReply{Message: "Heloo" + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb_test.RegisterGreeterServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
