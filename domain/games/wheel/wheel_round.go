package wheel

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	DefaultWheelRoundStartDelay = 10 * time.Second
	DefaultWheelRoundDuration   = 10 * time.Second
)

// TODO: Wheel -> Wheel/Fortune Wheel
type WheelRound struct {
	ProvablyFair
	lock      *sync.Mutex
	ID        uuid.UUID
	Entries   []WheelRoundEntry
	Outcome   wheelItem
	StartTime time.Time
	EndTime   time.Time
}

type WheelRoundEntry struct {
	UserID string
	Wager  float64
	Pick   WheelItemColor
}

func NewWheelRound(clientSeed, serverSeed []byte) (*WheelRound, error) {
	provablyFair, err := NewProvablyFair(clientSeed, serverSeed)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate a provably fair round")
	}

	startTime := time.Now().Add(DefaultWheelRoundStartDelay)
	endTime := startTime.Add(DefaultWheelRoundDuration)

	return &WheelRound{
		ProvablyFair: *provablyFair,
		ID:           uuid.New(),
		StartTime:    startTime,
		EndTime:      endTime,
	}, nil
}

func (r *WheelRound) Join(userId string, wager float64, pick WheelItemColor) {
	r.lock.Lock()
	defer r.lock.Unlock()

	entry := WheelRoundEntry{userId, wager, pick}

	r.Entries = append(r.Entries, entry)
}

func (r *WheelRound) Roll() (wheelItem, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	roll, err := r.Calculate()
	if err != nil {
		return wheelItem{}, errors.Wrap(err, "failed to calculate round roll")
	}

	r.Nonce++

	winningItem, err := GetItemByIdx(int(roll))
	if err != nil {
		return wheelItem{}, errors.Wrap(err, "failed to get winning item")
	}

	r.Outcome = winningItem

	return winningItem, nil
}

func Verify(clientSeed []byte, serverSeed []byte, nonce uint64, randNum uint64) (bool, error) {
	game, _ := NewWheelRound(clientSeed, serverSeed)
	game.Nonce = nonce

	roll, err := game.Calculate()

	log.Println(roll)

	if err != nil {
		return false, errors.Wrap(err, "failed to verify round")
	}

	return roll == randNum, nil
}

func (r *WheelRound) GetOutcomeIdx() int {
	return r.Outcome.idx
}
