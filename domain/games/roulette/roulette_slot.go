package roulette

import "github.com/pkg/errors"

var (
	errorIdxOutOfBounds = errors.New("slot index out of valid bounds (0 - 14)")
)

type RouletteSlotColor string

const (
	SlotColorRed   RouletteSlotColor = "red"
	SlotColorGreen RouletteSlotColor = "green"
	SlotColorBlack RouletteSlotColor = "black"
)

type rouletteSlot struct {
	idx    int
	color  RouletteSlotColor
	payout float64
}

const (
	WinningPercentageRed   = float64(14) / float64(15)
	WinningPercentageBlack = float64(14) / float64(15)
	WinningPercentageGreen = float64(1) / float64(15)
)

var (
	rouletteSlots = []rouletteSlot{
		{idx: 0, color: SlotColorBlack, payout: 1},
		{idx: 1, color: SlotColorRed, payout: 1},
		{idx: 2, color: SlotColorBlack, payout: 1},
		{idx: 3, color: SlotColorRed, payout: 1},
		{idx: 4, color: SlotColorBlack, payout: 1},
		{idx: 5, color: SlotColorRed, payout: 1},
		{idx: 6, color: SlotColorBlack, payout: 1},

		{idx: 7, color: SlotColorGreen, payout: 14},

		{idx: 8, color: SlotColorRed, payout: 1},
		{idx: 9, color: SlotColorBlack, payout: 1},
		{idx: 10, color: SlotColorRed, payout: 1},
		{idx: 11, color: SlotColorBlack, payout: 1},
		{idx: 12, color: SlotColorRed, payout: 1},
		{idx: 13, color: SlotColorBlack, payout: 1},
		{idx: 14, color: SlotColorRed, payout: 1},
	}
)

func GetSlotByIdx(idx uint64) (rouletteSlot, error) {
	if idx > 14 {
		return rouletteSlot{}, errorIdxOutOfBounds
	}

	return rouletteSlots[idx], nil
}

func GetRouletteSlots() []rouletteSlot {
	return rouletteSlots
}
