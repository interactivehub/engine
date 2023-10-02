package user

import "errors"

type User struct {
	ID       string
	UniqueID string
	Nickname string
	Points   int
}

func NewUser(id, uniqueId, nickname string, points int) (*User, error) {
	if id == "" {
		return nil, errors.New("missing id")
	}

	if uniqueId == "" {
		return nil, errors.New("missing uniqueId")
	}

	if nickname == "" {
		return nil, errors.New("missing nickname")
	}

	return &User{
		id,
		uniqueId,
		nickname,
		points,
	}, nil
}

func (u User) CanBet(bet int) bool {
	return u.Points >= bet
}
