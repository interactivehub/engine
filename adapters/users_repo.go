package adapters

import (
	"context"

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
	id       string `db:"id"`
	uniqueId string `db:unique_id`
	nickname string `db:nickname`
	points   int    `db:points`
}

func NewFromUser(user user.User) *sqlUser {
	return &sqlUser{
		id:       user.ID(),
		uniqueId: user.UniqueID(),
		nickname: user.Nickname(),
		points:   user.Points(),
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

	err := u.db.Get(&user, GetUserByIDQuery, id)
	if err != nil {
		return user, errors.Wrap(err, "failed to get user by id")
	}

	return user, nil
}

func (u UsersRepo) UserWithIdExists(ctx context.Context, id string) (bool, error) {
	count := 0

	err := u.db.Get(&count, CountUserByIdQuery, id)
	if err != nil {
		return false, errors.Wrap(err, "failed to count users by id")
	}

	return count > 0, nil
}

func (u UsersRepo) CreateUser(ctx context.Context, user user.User) error {
	sqlUser := NewFromUser(user)

	rows, err := u.db.NamedQuery(CreateUserQuery, sqlUser) // TODO fix this shit
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	defer rows.Close()

	return nil
}
