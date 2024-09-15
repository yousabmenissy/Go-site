package users

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID         int64
	Name       string
	Email      string
	Password   password
	Activated  bool
	Created_at time.Time
}

type password struct {
	plaintext *string
	Hash      []byte
}

func (p *password) Set(passwordPlainText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwordPlainText), 12)
	if err != nil {
		return err
	}

	p.plaintext = &passwordPlainText
	p.Hash = hash

	return nil
}

func (p *password) Matches(passwordPlainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(passwordPlainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil

		default:
			return false, err
		}
	}

	return true, nil
}

func (u *UserModel) Exists(id int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT true FROM users WHERE id = $1)`

	err := u.DB.QueryRow(query, id).Scan(&exists)
	return exists, err
}
