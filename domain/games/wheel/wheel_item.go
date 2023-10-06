package wheel

import "github.com/pkg/errors"

var (
	errorIdxOutOfBounds = errors.New("item index out of valid bounds (0 - 14)")
)

type WheelItemColor string

const (
	ItemColorGrey   WheelItemColor = "grey"
	ItemColorBlue   WheelItemColor = "blue"
	ItemColorYellow WheelItemColor = "yellow"
	ItemColorBrand  WheelItemColor = "brand"
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
	WinningPercentageBrand  = float64(1) / float64(31)
)

var (
	wheelItems = []wheelItem{
		{color: "#FA4475", payout: WinningPercentageBrand},  // Brand, 30x
		{color: "#FFE066", payout: WinningPercentageYellow}, // Red, 5x
		{color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{color: "#5ADBFF", payout: WinningPercentageBlue},   // Blue, 3x
		{color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{color: "#5ADBFF", payout: WinningPercentageBlue},   // Blue, 3x
		{color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{color: "#5ADBFF", payout: WinningPercentageBlue},   // Blue, 3x
		{color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{color: "#FFE066", payout: WinningPercentageYellow}, // Red, 5x
		{color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{color: "#FFE066", payout: WinningPercentageYellow}, // Red, 5x
		{color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{color: "#5ADBFF", payout: WinningPercentageBlue},   // Blue, 3x
		{color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{color: "#5ADBFF", payout: WinningPercentageBlue},   // Blue, 3x
		{color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{color: "#5ADBFF", payout: WinningPercentageBlue},   // Blue, 3x
		{color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{color: "#5ADBFF", payout: WinningPercentageBlue},   // Blue, 3x
		{color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{color: "#FFE066", payout: WinningPercentageYellow}, // Red, 5x
		{color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{color: "#FFE066", payout: WinningPercentageYellow}, // Red, 5x
		{color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{color: "#5ADBFF", payout: WinningPercentageBlue},   // Blue, 3x
		{color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{color: "#5ADBFF", payout: WinningPercentageBlue},   // Blue, 3x
		{color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
		{color: "#5ADBFF", payout: WinningPercentageBlue},   // Blue, 3x
		{color: "#606E80", payout: WinningPercentageGrey},   // Grey, 2x
	}
)

func GetItemByIdx(idx int) (wheelItem, error) {
	if idx > len(wheelItems)-1 {
		return wheelItem{}, errorIdxOutOfBounds
	}

	return wheelItems[idx], nil
}

func GetWheelItems() []wheelItem {
	return wheelItems
}
