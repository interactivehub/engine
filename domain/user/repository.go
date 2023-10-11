package user

import (
	"context"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	UserWithIdExists(ctx context.Context, id string) (bool, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
}
