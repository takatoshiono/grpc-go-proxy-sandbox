package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/utahta/grpc-go-proxy-example/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	address     = "localhost:50052" // server:50051, proxy:50052
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := helloworld.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Outgoing: abe hiroshi")
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("abe", "hiroshi"))

	r, err := c.SayHello(ctx, &helloworld.HelloRequest{Name: name})
	if err != nil {
		log.Printf("could not greet: %v\n", err)
		s, ok := status.FromError(err)
		if ok {
			details := s.Details()
			for i, detail := range details {
				switch d := detail.(type) {
				case *errdetails.BadRequest:
					log.Printf("BadRequest detail[%d]: %+v\n", i, d)
				default:
					log.Printf("Default detail[%d]: %+v\n", i, d)
				}
			}
		}
		return
	}
	log.Printf("Greeting: %s", r.Message)
}
