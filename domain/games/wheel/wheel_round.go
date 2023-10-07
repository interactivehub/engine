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
	ErrWheelRoundNotOpen    = errors.New("wheel round is not open")
	ErrWheelRoundNotRolling = errors.New("wheel round is not rolling")
)

type WheelRoundStatus string

const (
	WheelRoundStatusOpen    WheelRoundStatus = "open"
	WheelRoundStatusRolling WheelRoundStatus = "rolling"
	WheelRoundStatusClosed  WheelRoundStatus = "closed"
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

// TODO: Add More validations
func NewWheelRound(clientSeed, serverSeed []byte, prevNonce uint64) (*WheelRound, error) {
	provablyFair, err := NewProvablyFair(clientSeed, serverSeed, prevNonce)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate a provably fair round")
	}

	startTime := time.Now().Add(DefaultWheelRoundStartDelay)

	return &WheelRound{
		ProvablyFair: *provablyFair,
		ID:           uuid.New(),
		Status:       WheelRoundStatusOpen,
		StartTime:    startTime,
		lock:         new(sync.Mutex),
	}, nil
}

func (r *WheelRound) EndRound() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.Status != WheelRoundStatusRolling {
		return ErrWheelRoundNotRolling
	}

	endTime := time.Now()

	minRoundLength := DefaultWheelRoundStartDelay + DefaultWheelRoundDuration

	if endTime.Sub(r.StartTime) < minRoundLength {
		return ErrWheelRoundEndTooSoon
	}

	r.EndTime = time.Now()
	r.Status = WheelRoundStatusClosed

	return nil
}

func (r *WheelRound) Join(userId string, wager float64, pick WheelItemColor) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.Status != WheelRoundStatusOpen {
		return ErrWheelRoundNotOpen
	}

	entry := WheelRoundEntry{userId, wager, pick}

	r.Entries = append(r.Entries, entry)

	return nil
}

func (r *WheelRound) Roll() (wheelItem, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.Status != WheelRoundStatusOpen {
		return wheelItem{}, ErrWheelRoundNotOpen
	}

	roll, err := r.Calculate()
	if err != nil {
		return wheelItem{}, errors.Wrap(err, "failed to calculate round roll")
	}

	winningItem, err := GetItemByIdx(int(roll))
	if err != nil {
		return wheelItem{}, errors.Wrap(err, "failed to get winning item")
	}

	r.Outcome = winningItem
	r.Status = WheelRoundStatusRolling

	return winningItem, nil
}

func Verify(clientSeed []byte, serverSeed []byte, nonce uint64, randNum uint64) (bool, error) {
	game, _ := NewWheelRound(clientSeed, serverSeed, 0)
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

	return previousRound.Status == WheelRoundStatusClosed
}
