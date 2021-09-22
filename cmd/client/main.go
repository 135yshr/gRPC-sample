package main

import (
	"context"
	"fmt"
	"github.com/135yshr/gRPC_Sample/pkg/helloworld"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

func main() {
	fmt.Println("Run client...")

	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c:=helloworld.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	name := "World"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	r, err := c.SayHello(ctx, &helloworld.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Println("Greeting:", r.GetMessage())
}
