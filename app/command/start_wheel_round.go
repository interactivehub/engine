package command

import (
	"context"

	"github.com/interactivehub/engine/adapters"
	"github.com/interactivehub/engine/common/decorator"
	"github.com/interactivehub/engine/domain/games/wheel"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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
	logger *zap.Logger,
) StartWheelRoundHandler {
	if wsWriter == nil {
		panic("missing wsWriter")
	}

	if wheelRoundsRepo == nil {
		panic("missing wheelRoundsRepo")
	}

	return decorator.ApplyCommandDecorators[StartWheelRound](
		startWheelRoundHandler{wsWriter, wheelRoundsRepo},
		logger,
	)

}

func (h startWheelRoundHandler) Handle(ctx context.Context, cmd StartWheelRound) error {
	// 1. Retrieve latest wheel round
	latestRound, err := h.wheelRoundsRepo.GetLatest(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve latest wheel round")
	}

	// 2. Check if the latest round:
	//    - is nil: the current round is the first round
	//    - has status `end`: the latest round is ended
	canStart := wheel.CanStartNewRound(latestRound)
	if !canStart {
		return errors.New("cannot start wheel round")
	}

	// 3. Set the previous round nonce
	latestNonce := uint64(0)
	if latestRound != nil {
		latestNonce = latestRound.Nonce
	}

	// 4. Create a new wheel round
	round, err := wheel.NewWheelRound(cmd.ClientSeed, nil, latestNonce, wheel.WheelRoundOpenDuration, wheel.WheelRoundSpinDuration)
	if err != nil {
		return errors.Wrap(err, "failed to start wheel round")
	}

	// 4. Start an automatic wheel round handler
	round.
		Auto().
		Start().
		OnRoundStart(func(r *wheel.WheelRound) error {
			// TODO: notify client
			err := h.wheelRoundsRepo.Upsert(ctx, round)
			if err != nil {
				return errors.Wrap(err, "failed to upsert wheel round")
			}

			return nil
		}).
		OnStatusChange(func(r *wheel.WheelRound) error {
			// TODO: notify client
			return nil
		}).
		OnRoundEnd(func(r *wheel.WheelRound) error {
			err = h.wheelRoundsRepo.Upsert(ctx, round)
			if err != nil {
				return errors.Wrap(err, "failed to upsert wheel round")
			}

			return nil
		})

	return nil
}
