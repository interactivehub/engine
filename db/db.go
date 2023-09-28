package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func NewConnection() (*sqlx.DB, error) {
	dsn := getDSN()

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database connection")
	}

	return db, nil
}

func CloseConnection(db *sqlx.DB) error {
	err := db.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close database connection")
	}

	return nil
}

func getDSN() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host,
		port,
		name,
		user,
		password,
	)
}
