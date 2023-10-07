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
		return errors.Wrap(err, "failed to retrieve latest wheel round")
	}

	canStart := wheel.CanStartNewRound(latestRound)
	if !canStart {
		return ErrCannotStartWheelRound
	}

	round, err := wheel.NewWheelRound(cmd.ClientSeed, nil, latestRound.Nonce, wheel.WheelRoundOpenDuration, wheel.WheelRoundSpinDuration)
	if err != nil {
		return errors.Wrap(err, "failed to start wheel round")
	}

	round.
		Auto().
		Start().
		OnStatusChange(func(r *wheel.WheelRound) error {
			log.Println("status changed")
			log.Println(r.Status)

			return nil
		}).
		OnRoundEnd(func(r *wheel.WheelRound) error {
			err = h.wheelRoundsRepo.PersistWheelRound(ctx, *round)
			if err != nil {
				log.Println(err)
				return errors.Wrap(err, "failed to store wheel round")
			}

			return nil
		})

	return nil
}
