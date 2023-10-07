package app

import (
	"github.com/interactivehub/engine/app/command"
)

type Application struct {
	Commands Commands
}

type Commands struct {
	NewUser         command.NewUserHandler
	StartWheelRound command.StartWheelRoundHandler
}
