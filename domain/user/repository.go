package user

import (
	"context"
)

type Repository interface {
	GetUserById(ctx context.Context, id string) (User, error)
	UserWithIdExists(ctx context.Context, id string) (bool, error)
	CreateUser(ctx context.Context, user User) error
}
