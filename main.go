package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/interactivehub/engine/adapters"
	"github.com/interactivehub/engine/db"
	"github.com/interactivehub/engine/domain/user"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	db, err := db.NewConnection()
	if err != nil {
		panic(err)
	}

	usersRepo := adapters.NewUsersRepo(db)
	usersServer := adapters.NewUsersGrpcServer(usersRepo)

	user.RegisterUsersServiceServer(s, usersServer)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
