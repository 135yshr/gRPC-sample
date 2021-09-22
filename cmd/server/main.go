package main

import (
	"context"
	"fmt"
	"github.com/135yshr/gRPC_Sample/pkg/helloworld"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	fmt.Println("Run server...")
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, newServer())
	if err := s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

type server struct {
	helloworld.UnimplementedGreeterServer
}

func (h server) SayHello(_ context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{
		Message: fmt.Sprintf("Hello %s", req.GetName()),
	}, nil
}

func newServer() *server {
	s := &server{}
	return s
}
