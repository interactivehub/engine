package user

import (
	"errors"

	"github.com/interactivehub/engine/domain/currency"
)

type User struct {
	ID       string
	UniqueID string
	Nickname string
	HubMoney *currency.HubMoney
}

func NewUser(id, uniqueId, nickname string, amount float64) (*User, error) {
	if id == "" {
		return nil, errors.New("missing id")
	}

	if uniqueId == "" {
		return nil, errors.New("missing uniqueId")
	}

	if nickname == "" {
		return nil, errors.New("missing nickname")
	}

	hubMoney := currency.NewHubMoney(amount)

	if hubMoney.IsNegative() {
		return nil, errors.New("negative hub money")
	}

	return &User{
		ID:       id,
		UniqueID: uniqueId,
		Nickname: nickname,
		HubMoney: hubMoney,
	}, nil
}

func (u User) HasEnoughBalance(bet float64) bool {
	hubBet := currency.NewHubMoney(bet)

	hasEnough, err := u.HubMoney.GreaterThanOrEqual(hubBet.Money)
	if err != nil {
		return false
	}

	return hasEnough
}
