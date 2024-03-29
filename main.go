package main

import (
	"os"
	"strconv"

	"github.com/interactivehub/engine/adapters"
	"github.com/interactivehub/engine/app"
	"github.com/interactivehub/engine/app/command"
	"github.com/interactivehub/engine/common/db"
	"github.com/interactivehub/engine/common/logger"
	"github.com/interactivehub/engine/common/server"
	"github.com/interactivehub/engine/domain/games/wheel"
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

	logger := logger.Init()
	defer logger.Sync()

	// TODO: Find a way to set this bitchass inside app.NewApplication
	usersRepo := adapters.NewUsersRepo(db)
	wheelRoundsRepo := adapters.NewWheelRoundsRepo(db)
	wsWriter := adapters.NewWSWriter()

	app := app.Application{
		Commands: app.Commands{
			NewUser:         command.NewNewUserHandler(usersRepo, wsWriter, logger),
			StartWheelRound: command.NewStartWheelRoundHandler(wsWriter, wheelRoundsRepo, logger),
			JoinWheelRound:  command.NewJoinWheelRoundHandler(wsWriter, wheelRoundsRepo, usersRepo, logger),
		},
	}

	if runWS {
		go server.RunWSServer(func(server *server.WSServer) {
			wsWriter.SetClient(server.Client())

			wsListener := ports.NewWSListener(server.Client(), app, logger)
			wsListener.ListenEvents()
		})
	}

	if runGRPC {
		go server.RunGRPCServer(func(server *grpc.Server) {
			usersGrpcServer := ports.NewUsersGrpcServer(app)
			user.RegisterUsersServiceServer(server, usersGrpcServer)

			wheelGrpcService := ports.NewWheelGrpcServer(app)
			wheel.RegisterWheelServiceServer(server, wheelGrpcService)
		})
	}

	select {}
}
