package adapters

import (
	"context"
	"database/sql"

	"github.com/interactivehub/engine/domain/user"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	CreateUserQuery    = "INSERT INTO users (id, unique_id, nickname, points) VALUES (:id, :unique_id, :nickname, :points)"
	CountUserByIdQuery = "SELECT COUNT(id) FROM users WHERE id=$1"
	GetUserByIDQuery   = "SELECT id, unique_id, nickname, points FROM users WHERE id=$1"
)

type sqlUser struct {
	ID       string `db:"id"`
	UniqueID string `db:"unique_id"`
	Nickname string `db:"nickname"`
	Points   int    `db:"points"`
}

type UsersRepo struct {
	db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) *UsersRepo {
	if db == nil {
		panic("missing db")
	}

	return &UsersRepo{db}
}

func (u UsersRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	sqlUser := &sqlUser{}

	if err := u.db.GetContext(ctx, sqlUser, GetUserByIDQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to get user by id")
	}

	return sqlUser.toUser(), nil
}

func (u UsersRepo) UserWithIdExists(ctx context.Context, id string) (bool, error) {
	count := 0

	err := u.db.GetContext(ctx, &count, CountUserByIdQuery, id)
	if err != nil {
		return false, errors.Wrap(err, "failed to count users by id")
	}

	return count > 0, nil
}

func (u UsersRepo) CreateUser(ctx context.Context, user *user.User) error {
	sqlUser := &sqlUser{}
	sqlUser.fromUser(user)

	rows, err := u.db.NamedQueryContext(ctx, CreateUserQuery, sqlUser)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	defer rows.Close()

	return nil
}

func (u *sqlUser) fromUser(user *user.User) {
	if user == nil {
		return
	}

	u.ID = user.ID
	u.UniqueID = user.UniqueID
	u.Nickname = user.Nickname
	u.Points = user.Points
}

func (u *sqlUser) toUser() *user.User {
	if u == nil {
		return nil
	}

	return &user.User{
		ID:       u.ID,
		UniqueID: u.UniqueID,
		Nickname: u.Nickname,
		Points:   u.Points,
	}
}
