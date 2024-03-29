package adapters

import (
	"context"
	"database/sql"

	"github.com/interactivehub/engine/domain/user"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	CreateUserQuery = "INSERT INTO users (id, unique_id, nickname, hub_money) VALUES (:id, :unique_id, :nickname, :hub_money)"
	UpdateUserQuery = `
            UPDATE users 
            SET 
                id = :id,
                unique_id = :unique_id,
                nickname = :nickname,
                hub_money = :hub_money,
            WHERE id = :id
            `
	CountUserByIdQuery = "SELECT COUNT(id) FROM users WHERE id=$1"
	GetUserByIDQuery   = "SELECT id, unique_id, nickname, hub_money FROM users WHERE id=$1"
)

type sqlUser struct {
	ID       string  `db:"id"`
	UniqueID string  `db:"unique_id"`
	Nickname string  `db:"nickname"`
	HubMoney float64 `db:"hub_money"`
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

	_, err := u.db.NamedExecContext(ctx, CreateUserQuery, sqlUser)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	return nil
}

func (u UsersRepo) UpdateUser(ctx context.Context, user *user.User) error {
	userExists, err := u.UserWithIdExists(ctx, user.ID)
	if err != nil {
		return errors.Wrap(err, "failed to get user by id")
	}

	if !userExists {
		return errors.New("failed to update user: unknown id")
	}

	sqlUser := &sqlUser{}
	sqlUser.fromUser(user)

	_, err = u.db.NamedExecContext(ctx, UpdateUserQuery, sqlUser)
	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	return nil
}

func (u *sqlUser) fromUser(user *user.User) {
	if user == nil {
		return
	}

	u.ID = user.ID
	u.UniqueID = user.UniqueID
	u.Nickname = user.Nickname
	u.HubMoney = user.HubMoney.AsMajorUnits()
}

func (u *sqlUser) toUser() *user.User {
	if u == nil {
		return nil
	}

	user, err := user.NewUser(u.ID, u.UniqueID, u.Nickname, u.HubMoney)
	if err != nil {
		panic(err)
	}

	return user
}
