package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/utahta/grpc-go-proxy-example/helloworld"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	if t, ok := ctx.Deadline(); ok {
		fmt.Printf("in DEADLINE: %v\n", t)
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Printf("Incoming: abe %v\n", md.Get("abe")[0])
	}

	log.Printf("Received: %v\n", in.Name)

	if in.Name == "error" {
		s := status.New(codes.InvalidArgument, "invalid argument error")
		s, err := s.WithDetails(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "Name",
					Description: "error desu",
				},
			},
		})
		if err != nil {
			s = status.New(codes.Internal, fmt.Sprintf("failed to append invalid argument details message: %v", err))
		}
		return nil, s.Err()
	}

	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
