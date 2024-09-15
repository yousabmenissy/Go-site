package users

import (
	"database/sql"
	"errors"
	"site/internal"
)

func (u *UserModel) UpdateUser(user *User) error {
	query := `UPDATE users SET ("name", "email", "password_hash", "activated") = ($2, $3, $4, $5) WHERE id = $1`
	_, err := u.DB.Exec(query, user.ID, user.Name, user.Email, user.Password.Hash, user.Activated)
	if err != nil {

		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return internal.ErrDuplicateEmail

		case errors.Is(err, sql.ErrNoRows):
			return internal.ErrEditConflict

		default:
			return err
		}
	}
	return nil
}
