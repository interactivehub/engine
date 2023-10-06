package wheel

import (
	"context"
)

type Repository interface {
	PersistWheelRound(ctx context.Context, round WheelRound) error
	GetLatest(ctx context.Context) (WheelRound, error)
}
