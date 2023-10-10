package user

import "errors"

type User struct {
	ID       string
	UniqueID string
	Nickname string
	Points   float64
}

func NewUser(id, uniqueId, nickname string, points float64) (*User, error) {
	if id == "" {
		return nil, errors.New("missing id")
	}

	if uniqueId == "" {
		return nil, errors.New("missing uniqueId")
	}

	if nickname == "" {
		return nil, errors.New("missing nickname")
	}

	if points < 0 {
		return nil, errors.New("negative points")
	}

	return &User{
		id,
		uniqueId,
		nickname,
		points,
	}, nil
}

func (u User) HasEnoughBalance(bet float64) bool {
	return u.Points >= bet
}
