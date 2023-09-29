package server

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func RunGRPCServer(registerServer func(server *grpc.Server)) {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	addr := fmt.Sprintf(":%s", port)
	RunGRPCServerOnAddr(addr, registerServer)
}

func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
	grpcServer := grpc.NewServer()
	registerServer(grpcServer)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(grpcServer.Serve(listen))
}
