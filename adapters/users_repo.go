package adapters

import (
	"context"

	"github.com/interactivehub/engine/domain/user"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	CreateUserQuery     = "INSERT INTO users (id, unique_id, nickname, points) VALUES (:id, :unique_id, :nickname, :points)"
	CountUserByIdQuery  = "SELECT COUNT(id) FROM users WHERE id=$1"
	GetUserByIDQuery    = "SELECT id, unique_id, nickname, points FROM users WHERE id=$1"
	GetLeaderBoardQuery = "SELECT id, unique_id, nickname, points FROM users ORDER BY points DESC LIMIT $1"
)

type sqlUser struct {
	ID       string `db:"id"`
	UniqueID string `db:"unique_id"`
	Nickname string `db:"nickname"`
	Points   int    `db:"points"`
}

func newFromUser(user user.User) *sqlUser {
	return &sqlUser{
		ID:       user.ID(),
		UniqueID: user.UniqueID(),
		Nickname: user.Nickname(),
		Points:   user.Points(),
	}
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

func (u UsersRepo) TableName() string {
	return "users"
}

func (u UsersRepo) GetUserById(ctx context.Context, id string) (user.User, error) {
	user := user.User{}

	err := u.db.GetContext(ctx, &user, GetUserByIDQuery, id)
	if err != nil {
		return user, errors.Wrap(err, "failed to get user by id")
	}

	return user, nil
}

func (u UsersRepo) GetLeaderBoard(ctx context.Context, limit int32) ([]user.User, error) {
	var leaderBoard []user.User

	err := u.db.SelectContext(ctx, &leaderBoard, GetLeaderBoardQuery, limit)
	if err != nil {
		return leaderBoard, errors.Wrap(err, "failed to get leader board")
	}

	return leaderBoard, nil
}

func (u UsersRepo) UserWithIdExists(ctx context.Context, id string) (bool, error) {
	count := 0

	err := u.db.GetContext(ctx, &count, CountUserByIdQuery, id)
	if err != nil {
		return false, errors.Wrap(err, "failed to count users by id")
	}

	return count > 0, nil
}

func (u UsersRepo) CreateUser(ctx context.Context, user user.User) error {
	sqlUser := newFromUser(user)

	rows, err := u.db.NamedQueryContext(ctx, CreateUserQuery, sqlUser)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	defer rows.Close()

	return nil
}
