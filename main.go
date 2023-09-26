package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/interactivehub/engine/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedUsersServiceServer
}

func (s *server) NewUser(ctx context.Context, in *pb.NewUserRequest) (*pb.NewUserResponse, error) {
	log.Printf("Received: %v", in.GetUserId())
	return &pb.NewUserResponse{
		Success: true,
		Data: "Hello",
	}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUsersServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}