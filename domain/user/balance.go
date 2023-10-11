package user

import (
	"errors"

	"github.com/interactivehub/engine/domain/currency"
)

func (u *User) HasHubMoney(money float64) bool {
	hubMoney := currency.NewHubMoney(money)

	hasEnough, err := u.HubMoney.GreaterThanOrEqual(hubMoney.Money)
	if err != nil {
		return false
	}

	return hasEnough
}

func (u *User) Bet(bet float64) error {
	hubBet := currency.NewHubMoney(bet)

	if hubBet.IsNegative() {
		return errors.New("negative bet")
	}

	if !u.HasHubMoney(bet) {
		return errors.New("user has not enough balance")
	}

	_, err := u.HubMoney.Subtract(hubBet.Money)
	return err
}
