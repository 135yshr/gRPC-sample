package main

import (
	"context"
	"fmt"
	"github.com/135yshr/gRPC_Sample/pkg/helloworld"
	"github.com/135yshr/gRPC_Sample/pkg/knockknock"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"time"
)

var name = "World"

func main() {
	fmt.Println("Run client...")
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	gc := helloworld.NewGreeterClient(conn)
	if err := sendGreeter(gc); err != nil {
		log.Fatalf("Failed send greeter: %v", err)
	}

	waitc := make(chan struct{})
	kc := knockknock.NewDoorClient(conn)
	go func() {
		err := sendKnock(kc)
		if err != nil {
			log.Fatalf("Failed send greeter: %v", err)
		}
		close(waitc)
	}()

	if err := sendGreeter(gc); err != nil {
		log.Fatalf("Failed send greeter: %v", err)
	}
	time.Sleep(3 * time.Second)
	if err := sendGreeter(gc); err != nil {
		log.Fatalf("Failed send greeter: %v", err)
	}
	<-waitc
}

func receive(stream knockknock.Door_KnockClient) chan error {
	cerr := make(chan error)
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				log.Println("EOF")
				close(cerr)
				return
			}
			if err != nil {
				cerr <- err
			}
			fmt.Println(in.Message)
		}
	}()
	return cerr
}
func sendKnock(kc knockknock.DoorClient) error {
	stream, err := kc.Knock(context.Background())
	if err != nil {
		return err
	}
	req := &knockknock.KnockRequest{Count: 10, Wait: 1}
	if err := stream.Send(req); err != nil {
		return err
	}
	log.Println("Send knock")
	select {
	case e := <-receive(stream):
		return e
	}
}

func sendGreeter(gc helloworld.GreeterClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := gc.SayHello(ctx, &helloworld.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Println("Greeting:", r.GetMessage())
	return nil
}
