package wheel

import (
	"context"
)

type Repository interface {
	PersistWheelRound(ctx context.Context, round WheelRound) error
}
