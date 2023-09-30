package main

import (
	"os"
	"strconv"

	"github.com/interactivehub/engine/adapters"
	"github.com/interactivehub/engine/app"
	"github.com/interactivehub/engine/app/handlers"
	"github.com/interactivehub/engine/common/db"
	"github.com/interactivehub/engine/common/server"
	"github.com/interactivehub/engine/domain/user"
	"github.com/interactivehub/engine/ports"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	runGRPC, _ := strconv.ParseBool(os.Getenv("RUN_GRPC"))
	runWS, _ := strconv.ParseBool(os.Getenv("RUN_WEBSOCKET"))

	db, err := db.NewConnection()
	if err != nil {
		panic(err)
	}

	// TODO: Find a way to set this bitchass inside app.NewApplication
	usersRepo := adapters.NewUsersRepo(db)
	wsWriter := adapters.NewWSWriter()

	app := app.Application{
		Handlers: app.Handlers{
			NewUser: handlers.NewNewUserHandler(usersRepo, wsWriter),
		},
	}

	if runGRPC {
		go server.RunGRPCServer(func(server *grpc.Server) {
			usersGrpcServer := ports.NewUsersGrpcServer(app)
			user.RegisterUsersServiceServer(server, usersGrpcServer)
		})
	}

	if runWS {
		go server.RunWSServer(func(server *server.WSServer) {
			wsListener := ports.NewWSListener(server.Client())
			wsWriter.SetClient(server.Client())
			wsListener.ListenEvents()
		})
	}

	select {}
}
