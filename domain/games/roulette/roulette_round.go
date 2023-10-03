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

// TODO: Persist this in db (for the future)
type RouletteRound struct {
	provablyFair
	id        uuid.UUID
	entries   []RouletteRoundEntry
	outcome   rouletteSlot
	lock      sync.Mutex
	startTime time.Time
	endTime   time.Time
}

// TODO: Persist this in db (for the future)
type RouletteRoundEntry struct {
	userId string
	wager  float64
	pick   RouletteSlotColor
}

func NewRouletteRound(clientSeed, serverSeed []byte) (*RouletteRound, error) {
	provablyFair, err := NewProvablyFair(clientSeed, serverSeed)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate a provably fair round")
	}

	startTime := time.Now().Add(DefaultRouletteRoundStartDelay)
	endTime := startTime.Add(DefaultRouletteRoundDuration)

	return &RouletteRound{
		provablyFair: *provablyFair,
		id:           uuid.New(),
		startTime:    startTime,
		endTime:      endTime,
	}, nil
}

func (r *RouletteRound) Join(userId string, wager float64, pick RouletteSlotColor) {
	r.lock.Lock()
	defer r.lock.Unlock()

	entry := RouletteRoundEntry{userId, wager, pick}

	r.entries = append(r.entries, entry)
}

func (r *RouletteRound) Roll() (rouletteSlot, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	roll, err := r.Calculate()
	if err != nil {
		return rouletteSlot{}, errors.Wrap(err, "failed to calculate round roll")
	}

	r.nonce++

	winningSlot, err := GetSlotByIdx(roll)
	if err != nil {
		return rouletteSlot{}, errors.Wrap(err, "failed to get winning slot")
	}

	r.outcome = winningSlot

	return winningSlot, nil
}

func Verify(clientSeed []byte, serverSeed []byte, nonce uint64, randNum uint64) (bool, error) {
	game, _ := NewRouletteRound(clientSeed, serverSeed)
	game.nonce = nonce

	roll, err := game.Calculate()

	log.Println(roll)

	if err != nil {
		return false, errors.Wrap(err, "failed to verify round")
	}

	return roll == randNum, nil
}
