package command

import (
	"context"
	"log"

	"github.com/interactivehub/engine/adapters"
	"github.com/interactivehub/engine/common/decorator"
	"github.com/interactivehub/engine/domain/games/wheel"
	"github.com/pkg/errors"
)

var (
	ErrCannotStartWheelRound = errors.New("cannot start wheel round")
)

type StartWheelRound struct {
	ClientSeed []byte
}

type StartWheelRoundHandler decorator.CommandHandler[StartWheelRound]

type startWheelRoundHandler struct {
	wsWriter        adapters.WSWriter
	wheelRoundsRepo wheel.Repository
}

func NewStartWheelRoundHandler(
	wsWriter adapters.WSWriter,
	wheelRoundsRepo wheel.Repository,
) StartWheelRoundHandler {
	if wsWriter == nil {
		panic("missing wsWriter")
	}

	if wheelRoundsRepo == nil {
		panic("missing wheelRoundsRepo")
	}

	return startWheelRoundHandler{wsWriter, wheelRoundsRepo}
}

func (h startWheelRoundHandler) Handle(ctx context.Context, cmd StartWheelRound) error {
	latestRound, err := h.wheelRoundsRepo.GetLatest(ctx)
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "failed to retrieve latest wheel round")
	}

	canStart := wheel.CanStartNewRound(latestRound)
	if !canStart {
		return ErrCannotStartWheelRound
	}

	newRound, err := wheel.NewWheelRound(cmd.ClientSeed, nil, latestRound.Nonce)
	if err != nil {
		return errors.Wrap(err, "failed to start wheel round")
	}

	err = h.wheelRoundsRepo.PersistWheelRound(ctx, *newRound)
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "failed to store wheel round")
	}

	return nil
}
