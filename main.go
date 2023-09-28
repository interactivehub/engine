package main

import (
	"github.com/interactivehub/engine/adapters"
	"github.com/interactivehub/engine/db"
	"github.com/interactivehub/engine/domain/user"
	"github.com/interactivehub/engine/ports"
	"github.com/interactivehub/engine/server"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	db, err := db.NewConnection()
	if err != nil {
		panic(err)
	}

	usersRepo := adapters.NewUsersRepo(db)

	server.RunGRPCServer(func(server *grpc.Server) {
		usersGrpcServer := ports.NewUsersGrpcServer(usersRepo)
		user.RegisterUsersServiceServer(server, usersGrpcServer)
	})
}
