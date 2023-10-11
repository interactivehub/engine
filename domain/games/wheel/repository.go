package wheel

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*WheelRound, error)
	GetLatest(ctx context.Context) (*WheelRound, error)
	Upsert(ctx context.Context, round *WheelRound) error
	UpsertEntry(ctx context.Context, roundEntry *WheelRoundEntry) error
	GetEntry(ctx context.Context, roundID uuid.UUID, userID string, pick WheelItemColor) (*WheelRoundEntry, error)
}
