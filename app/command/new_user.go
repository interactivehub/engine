package command

import (
	"context"

	"github.com/interactivehub/engine/adapters"
	"github.com/interactivehub/engine/common/decorator"
	"github.com/interactivehub/engine/domain/user"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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

type NewUserHandler decorator.CommandHandler[NewUser]

type newUserHandler struct {
	usersRepo user.Repository
	wsWriter  adapters.WSWriter
}

func NewNewUserHandler(
	usersRepo user.Repository,
	wsWriter adapters.WSWriter,
	logger *zap.Logger,
) NewUserHandler {
	if usersRepo == nil {
		panic("missing usersRepo")
	}

	if wsWriter == nil {
		panic("missing wsWriter")
	}

	return decorator.ApplyCommandDecorators[NewUser](
		newUserHandler{usersRepo, wsWriter},
		logger,
	)

}

func (h newUserHandler) Handle(ctx context.Context, cmd NewUser) error {
	userExists, err := h.usersRepo.UserWithIdExists(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "failed to create new user")
	}

	if userExists {
		return errors.New("failed to create user. user already exists")
	}

	user, err := user.NewUser(cmd.ID, cmd.UniqueID, cmd.Nickname, 0)
	if err != nil {
		return errors.Wrap(err, "failed to create new user")
	}

	err = h.usersRepo.CreateUser(ctx, user)
	if err != nil {
		return errors.Wrap(err, "failed to create new user")
	}

	newUserEventPayload := NewUserEventPayload{
		ID:       user.ID,
		UniqueID: user.UniqueID,
		Nickname: user.Nickname,
		Points:   user.Points,
	}

	err = h.wsWriter.WriteEvent(NewUserEventType, newUserEventPayload)
	if err != nil {
		return errors.Wrap(err, "error sending newly created user event")
	}

	return nil
}
