package command

import (
	"context"

	"github.com/interactivehub/engine/adapters"
	"github.com/interactivehub/engine/common/decorator"
	"github.com/interactivehub/engine/domain/games/wheel"
	"github.com/interactivehub/engine/domain/user"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type JoinWheelRound struct {
	UserID string
	Bet    float64
	Pick   wheel.WheelItemColor
}

type JoinWheelRoundHandler decorator.CommandHandler[JoinWheelRound]

type joinWheelRoundHandler struct {
	wsWriter        adapters.WSWriter
	wheelRoundsRepo wheel.Repository
	usersRepo       user.Repository
}

func NewJoinWheelRoundHandler(
	wsWriter adapters.WSWriter,
	wheelRoundsRepo wheel.Repository,
	usersRepo user.Repository,
	logger *zap.Logger,
) JoinWheelRoundHandler {
	if wsWriter == nil {
		panic("missing wsWriter")
	}

	if wheelRoundsRepo == nil {
		panic("missing wheelRoundsRepo")
	}

	if usersRepo == nil {
		panic("missing usersRepo")
	}

	return decorator.ApplyCommandDecorators[JoinWheelRound](
		joinWheelRoundHandler{wsWriter, wheelRoundsRepo, usersRepo},
		logger,
	)
}

// TODO: Gotta test this
func (h joinWheelRoundHandler) Handle(ctx context.Context, cmd JoinWheelRound) error {
	// 1. Retrieve user by id
	user, err := h.usersRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve entry user")
	}

	// 2. Check if user exists
	if user == nil {
		return errors.New("failed to join wheel round: user id is unknown")
	}

	// 3. Retrieve latest wheel round
	latestRound, err := h.wheelRoundsRepo.GetLatest(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve latest wheel round")
	}

	// 4. Check if the latest round is open
	canJoin := latestRound.IsStatus(wheel.WheelRoundStatusOpen)
	if !canJoin {
		return errors.New("failed to join wheel round: round is not open")
	}

	// TODO
	// 5. Subtract user money and update db entry (gotta do it in the same transaction as point 6)
	user.Bet(cmd.Bet)

	// 6. Create new wheel round entry
	roundEntry := wheel.NewWheelRoundEntry(latestRound.ID, cmd.UserID, cmd.Bet, cmd.Pick)

	// 7. Persist new wheel round entry
	err = h.wheelRoundsRepo.UpsertEntry(ctx, roundEntry)
	if err != nil {
		return errors.Wrap(err, "failed to create wheel round")
	}

	return nil
}
