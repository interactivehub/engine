package handlers

import (
	"context"

	"github.com/interactivehub/engine/adapters"
	"github.com/interactivehub/engine/common/decorator"
	"github.com/interactivehub/engine/domain/user"
	"github.com/pkg/errors"
)

const (
	NewUserEventType = "newUser"
)

type NewUser struct {
	ID       string
	UniqueID string
	Nickname string
}

type NewUserEventPayload struct {
	ID       string `json:"id"`
	UniqueID string `json:"uniqueId"`
	Nickname string `json:"nickname"`
	Points   int    `json:"points"`
}

type NewUserResponse struct {
	ID       string
	UniqueID string
	Nickname string
	Points   int
}

type NewUserHandler decorator.Handler[NewUser, NewUserResponse]

type newUserHandler struct {
	usersRepo user.Repository
	wsWriter  adapters.WSWriter
}

func NewNewUserHandler(usersRepo user.Repository, wsWriter adapters.WSWriter) NewUserHandler {
	if usersRepo == nil {
		panic("missing usersRepo")
	}

	if wsWriter == nil {
		panic("missing wsWriter")
	}

	return newUserHandler{usersRepo, wsWriter}
}

func (h newUserHandler) Handle(ctx context.Context, cmd NewUser) (NewUserResponse, error) {
	userExists, err := h.usersRepo.UserWithIdExists(ctx, cmd.ID)
	if err != nil {
		return NewUserResponse{}, err
	}

	if userExists {
		return NewUserResponse{}, errors.New("failed to create user. user already exists")
	}

	user, err := user.NewUser(cmd.ID, cmd.UniqueID, cmd.Nickname, 0)
	if err != nil {
		return NewUserResponse{}, err
	}

	err = h.usersRepo.CreateUser(ctx, *user)
	if err != nil {
		return NewUserResponse{}, err
	}

	newUserEventPayload := NewUserEventPayload{
		ID:       user.ID(),
		UniqueID: user.UniqueID(),
		Nickname: user.Nickname(),
		Points:   user.Points(),
	}

	err = h.wsWriter.WriteEvent(NewUserEventType, newUserEventPayload)
	if err != nil {
		return NewUserResponse{}, errors.Wrap(err, "error sending newly created user")
	}

	return NewUserResponse{
		ID:       user.ID(),
		UniqueID: user.UniqueID(),
		Nickname: user.Nickname(),
		Points:   user.Points(),
	}, nil
}
