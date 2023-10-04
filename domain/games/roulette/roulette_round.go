package roulette

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	DefaultRouletteRoundStartDelay = 10 * time.Second
	DefaultRouletteRoundDuration   = 10 * time.Second
)

type RouletteRound struct {
	ProvablyFair
	lock      *sync.Mutex
	ID        uuid.UUID
	Entries   []RouletteRoundEntry
	Outcome   rouletteSlot
	StartTime time.Time
	EndTime   time.Time
}

type RouletteRoundEntry struct {
	UserID string
	Wager  float64
	Pick   RouletteSlotColor
}

func NewRouletteRound(clientSeed, serverSeed []byte) (*RouletteRound, error) {
	provablyFair, err := NewProvablyFair(clientSeed, serverSeed)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate a provably fair round")
	}

	startTime := time.Now().Add(DefaultRouletteRoundStartDelay)
	endTime := startTime.Add(DefaultRouletteRoundDuration)

	return &RouletteRound{
		ProvablyFair: *provablyFair,
		ID:           uuid.New(),
		StartTime:    startTime,
		EndTime:      endTime,
	}, nil
}

func (r *RouletteRound) Join(userId string, wager float64, pick RouletteSlotColor) {
	r.lock.Lock()
	defer r.lock.Unlock()

	entry := RouletteRoundEntry{userId, wager, pick}

	r.Entries = append(r.Entries, entry)
}

func (r *RouletteRound) Roll() (rouletteSlot, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	roll, err := r.Calculate()
	if err != nil {
		return rouletteSlot{}, errors.Wrap(err, "failed to calculate round roll")
	}

	r.Nonce++

	winningSlot, err := GetSlotByIdx(roll)
	if err != nil {
		return rouletteSlot{}, errors.Wrap(err, "failed to get winning slot")
	}

	r.Outcome = winningSlot

	return winningSlot, nil
}

func Verify(clientSeed []byte, serverSeed []byte, nonce uint64, randNum uint64) (bool, error) {
	game, _ := NewRouletteRound(clientSeed, serverSeed)
	game.Nonce = nonce

	roll, err := game.Calculate()

	log.Println(roll)

	if err != nil {
		return false, errors.Wrap(err, "failed to verify round")
	}

	return roll == randNum, nil
}

func (r *RouletteRound) GetOutcomeIdx() int {
	return r.Outcome.idx
}
