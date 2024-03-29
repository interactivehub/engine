package wheel

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	WheelRoundOpenDuration = 30 * time.Second
	WheelRoundSpinDuration = 10 * time.Second
)

type WheelRoundStatus string

const (
	WheelRoundStatusOpen WheelRoundStatus = "open"
	WheelRoundStatusSpin WheelRoundStatus = "spin"
	WheelRoundStatusEnd  WheelRoundStatus = "end"
)

type WheelRound struct {
	ProvablyFair
	lock           *sync.Mutex
	ID             uuid.UUID
	Status         WheelRoundStatus
	Entries        []*WheelRoundEntry
	Outcome        wheelItem
	RoundStartTime time.Time
	RoundEndTime   time.Time
	SpinStartTime  time.Time
	OpenDuration   time.Duration
	SpinDuration   time.Duration
}

type WheelRoundEntry struct {
	RoundID   uuid.UUID
	UserID    string
	Bet       float64
	Pick      WheelItemColor
	EntryTime time.Time
}

func NewWheelRound(
	clientSeed,
	serverSeed []byte,
	prevNonce uint64,
	openDuration,
	spinDuration time.Duration,
) (*WheelRound, error) {
	provablyFair, err := NewProvablyFair(clientSeed, serverSeed, prevNonce)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate a provably fair round")
	}

	roundStartTime := time.Now()
	spinStartTime := roundStartTime.Add(openDuration)
	roundEndTime := spinStartTime.Add(spinDuration)

	return &WheelRound{
		ProvablyFair:   *provablyFair,
		ID:             uuid.New(),
		Status:         WheelRoundStatusOpen,
		lock:           new(sync.Mutex),
		RoundStartTime: roundStartTime,
		SpinStartTime:  spinStartTime,
		RoundEndTime:   roundEndTime,
		OpenDuration:   openDuration,
		SpinDuration:   spinDuration,
	}, nil
}

func NewWheelRoundEntry(
	roundId uuid.UUID,
	userId string,
	bet float64,
	pick WheelItemColor,
) *WheelRoundEntry {
	entryTime := time.Now()

	return &WheelRoundEntry{
		RoundID:   roundId,
		UserID:    userId,
		Bet:       bet,
		Pick:      pick,
		EntryTime: entryTime,
	}
}

func (r *WheelRound) Join(userId string, bet float64, pick WheelItemColor) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.Status != WheelRoundStatusOpen {
		return errors.New("wheel round is not open")
	}

	entry := NewWheelRoundEntry(r.ID, userId, bet, pick)

	r.Entries = append(r.Entries, entry)

	return nil
}

func (r *WheelRound) Roll() (wheelItem, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.Status != WheelRoundStatusOpen {
		return wheelItem{}, errors.New("wheel round is not open")
	}

	now := time.Now()
	if now.Before(r.SpinStartTime) || now.After(r.RoundEndTime) {
		return wheelItem{}, errors.New("rolling outside valid time")
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
	r.Status = WheelRoundStatusSpin

	return winningItem, nil
}

func (r *WheelRound) EndRound() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.Status != WheelRoundStatusSpin {
		return errors.New("wheel round is not in the spin phase")
	}

	now := time.Now()
	if now.Before(r.RoundEndTime) {
		return errors.New("ending the round too early")
	}

	r.Status = WheelRoundStatusEnd

	return nil
}

func (r *WheelRound) GetOutcomeIdx() int {
	return r.Outcome.idx
}

func (r *WheelRound) IsStatus(status WheelRoundStatus) bool {
	if r == nil {
		return false
	}

	return r.Status == status
}

func Verify(clientSeed []byte, serverSeed []byte, nonce uint64, randNum uint64) (bool, error) {
	game, _ := NewWheelRound(clientSeed, serverSeed, 0, WheelRoundOpenDuration, WheelRoundSpinDuration)
	game.Nonce = nonce

	roll, err := game.Calculate()

	if err != nil {
		return false, errors.Wrap(err, "failed to verify round")
	}

	return roll == randNum, nil
}

func CanStartNewRound(previousRound *WheelRound) bool {
	if previousRound == nil {
		return true
	}

	return previousRound.IsStatus(WheelRoundStatusEnd)
}
