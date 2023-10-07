package wheel

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	DefaultWheelRoundStartDelay = 15 * time.Second
	DefaultWheelRoundDuration   = 15 * time.Second
)

var (
	ErrWheelRoundEndTooSoon = errors.New("could not end wheel round yet")
	ErrWheelRoundNotOngoing = errors.New("wheel round is not ongoing")
)

type WheelRoundStatus string

const (
	WheelRoundOngoing WheelRoundStatus = "ongoing"
	WheelRoundEnded   WheelRoundStatus = "ended"
)

type WheelRound struct {
	ProvablyFair
	lock      *sync.Mutex
	ID        uuid.UUID
	Status    WheelRoundStatus
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

	return &WheelRound{
		ProvablyFair: *provablyFair,
		ID:           uuid.New(),
		Status:       WheelRoundOngoing,
		StartTime:    startTime,
		lock:         new(sync.Mutex),
	}, nil
}

func (r *WheelRound) EndRound() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.Status != WheelRoundOngoing {
		return ErrWheelRoundNotOngoing
	}

	endTime := time.Now()

	minRoundLength := DefaultWheelRoundStartDelay + DefaultWheelRoundDuration

	if endTime.Sub(r.StartTime) < minRoundLength {
		return ErrWheelRoundEndTooSoon
	}

	r.EndTime = time.Now()
	r.Status = WheelRoundEnded

	return nil
}

func (r *WheelRound) Join(userId string, wager float64, pick WheelItemColor) {
	r.lock.Lock()
	defer r.lock.Unlock()

	// TODO: Validate

	entry := WheelRoundEntry{userId, wager, pick}

	r.Entries = append(r.Entries, entry)
}

func (r *WheelRound) Roll() (wheelItem, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	// TODO: Check if round is ongoing

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

	if err != nil {
		return false, errors.Wrap(err, "failed to verify round")
	}

	return roll == randNum, nil
}

func (r *WheelRound) GetOutcomeIdx() int {
	return r.Outcome.idx
}

func CanStartNewRound(previousRound WheelRound) bool {
	if previousRound.ID == uuid.Nil {
		return true
	}

	return previousRound.Status == WheelRoundEnded
}
