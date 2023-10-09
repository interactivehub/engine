package wheel

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*WheelRound, error)
	GetLatest(ctx context.Context) (*WheelRound, error)
	Create(ctx context.Context, round *WheelRound) error
	CreateEntry(ctx context.Context, roundEntry *WheelRoundEntry) error
}
