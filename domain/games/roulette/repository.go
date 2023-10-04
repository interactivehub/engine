package roulette

import (
	"context"
)

type Repository interface {
	PersistRouletteRound(ctx context.Context, round RouletteRound) error
}
