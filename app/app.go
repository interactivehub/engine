package app

import (
	"github.com/interactivehub/engine/app/handlers"
)

type Application struct {
	Handlers Handlers
}

type Handlers struct {
	NewUser handlers.NewUserHandler
}
