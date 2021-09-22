package main

import (
	"context"
	"fmt"
	"github.com/135yshr/gRPC_Sample/pkg/helloworld"
	"github.com/135yshr/gRPC_Sample/pkg/knockknock"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	fmt.Println("Run server...")
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	s := grpc.NewServer()
	ss := &server{}
	helloworld.RegisterGreeterServer(s, ss)
	knockknock.RegisterDoorServer(s, ss)
	if err := s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

type server struct {
	helloworld.UnimplementedGreeterServer
	knockknock.UnimplementedDoorServer
}

func (s *server) SayHello(_ context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{
		Message: fmt.Sprintf("Hello %s", req.GetName()),
	}, nil
}

func (s *server) Knock(stream knockknock.Door_KnockServer) error {
	for {
		log.Println("Wait recv")
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Println("Receive request")
		count := in.Count
		if count == 0 {
			count = 100
		}
		wait := time.Duration(in.Wait) * time.Second
		if wait == 0 {
			wait = time.Second
		}

		select {
		case err := <-send(stream, int(count), wait):
			return err
		default:
			time.Sleep(time.Second)
		}
	}
}

func send(stream knockknock.Door_KnockServer, count int, wait time.Duration) chan error {
	cerr := make(chan error)
	go func() {
		for i := 0; i < count; i++ {
			if err := stream.Send(&knockknock.KnockReply{Message: "knock!"}); err != nil {
				cerr <- err
			}
			time.Sleep(wait)
		}
		close(cerr)
	}()
	return cerr
}
