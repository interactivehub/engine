package wheel

import (
	"strings"

	"github.com/pkg/errors"
)

var (
	errorIdxOutOfBounds = errors.New("item index out of valid bounds (0 - 30)")
)

type WheelItemColor string

const (
	ItemColorGrey   WheelItemColor = "grey"
	ItemColorBlue   WheelItemColor = "blue"
	ItemColorYellow WheelItemColor = "yellow"
	ItemColorRed    WheelItemColor = "red"
)

type wheelItem struct {
	idx    int
	color  WheelItemColor
	payout float64
}

const (
	WinningPercentageGrey   = float64(15) / float64(31)
	WinningPercentageBlue   = float64(10) / float64(31)
	WinningPercentageYellow = float64(5) / float64(31)
	WinningPercentageRed    = float64(1) / float64(31)
)

var (
	wheelItems = []wheelItem{
		{idx: 0, color: "#FA4475", payout: WinningPercentageRed},     // Red, 30x
		{idx: 1, color: "#FFE066", payout: WinningPercentageYellow},  // Red, 5x
		{idx: 2, color: "#606E80", payout: WinningPercentageGrey},    // Grey, 2x
		{idx: 3, color: "#5ADBFF", payout: WinningPercentageBlue},    // Blue, 3x
		{idx: 4, color: "#606E80", payout: WinningPercentageGrey},    // Grey, 2x
		{idx: 5, color: "#5ADBFF", payout: WinningPercentageBlue},    // Blue, 3x
		{idx: 6, color: "#606E80", payout: WinningPercentageGrey},    // Grey, 2x
		{idx: 7, color: "#5ADBFF", payout: WinningPercentageBlue},    // Blue, 3x
		{idx: 8, color: "#606E80", payout: WinningPercentageGrey},    // Grey, 2x
		{idx: 9, color: "#FFE066", payout: WinningPercentageYellow},  // Red, 5x
		{idx: 10, color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{idx: 11, color: "#FFE066", payout: WinningPercentageYellow}, // Red, 5x
		{idx: 12, color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{idx: 13, color: "#5ADBFF", payout: WinningPercentageBlue},   // Blue, 3x
		{idx: 14, color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{idx: 15, color: "#5ADBFF", payout: WinningPercentageBlue},   // Blue, 3x
		{idx: 16, color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{idx: 17, color: "#5ADBFF", payout: WinningPercentageBlue},   // Blue, 3x
		{idx: 18, color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{idx: 19, color: "#5ADBFF", payout: WinningPercentageBlue},   // Blue, 3x
		{idx: 20, color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{idx: 21, color: "#FFE066", payout: WinningPercentageYellow}, // Red, 5x
		{idx: 22, color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{idx: 23, color: "#FFE066", payout: WinningPercentageYellow}, // Red, 5x
		{idx: 24, color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{idx: 25, color: "#5ADBFF", payout: WinningPercentageBlue},   // Blue, 3x
		{idx: 26, color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{idx: 27, color: "#5ADBFF", payout: WinningPercentageBlue},   // Blue, 3x
		{idx: 28, color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{idx: 29, color: "#5ADBFF", payout: WinningPercentageBlue},   // Blue, 3x
		{idx: 30, color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
	}
)

func ParseWheelItemColor(s string) (WheelItemColor, error) {
	normalizedStr := strings.ToLower(s)

	colors := map[string]WheelItemColor{
		"grey":   ItemColorGrey,
		"blue":   ItemColorBlue,
		"yellow": ItemColorYellow,
		"red":    ItemColorRed,
	}

	itemColor, ok := colors[normalizedStr]
	if !ok {
		return itemColor, errors.New("failed to parse wheel item color")
	}

	return itemColor, nil
}

func WheelItems() []wheelItem {
	return wheelItems
}

func GetItemByIdx(idx int) (wheelItem, error) {
	if idx > len(wheelItems)-1 {
		return wheelItem{}, errorIdxOutOfBounds
	}

	return wheelItems[idx], nil
}

func GetWheelItems() []wheelItem {
	return wheelItems
}
