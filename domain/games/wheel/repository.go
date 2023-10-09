package wheel

import (
	"context"
)

type Repository interface {
	CreateWheelRound(ctx context.Context, round WheelRound) error
	GetLatest(ctx context.Context) (WheelRound, error)
}
