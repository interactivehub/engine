package user

import "errors"

type User struct {
	id       string
	uniqueId string
	nickname string
	points   int
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
	return u.points >= bet
}

func (u User) ID() string {
	return u.id
}

func (u User) UniqueID() string {
	return u.uniqueId
}

func (u User) Nickname() string {
	return u.nickname
}

func (u User) Points() int {
	return u.points
}
