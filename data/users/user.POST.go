package users

import "site/internal"

func (u *UserModel) Insert(user *User) error {
	statement := `INSERT INTO users ("name", "email", "password_hash", "activated") VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	if err := u.DB.QueryRow(statement, user.Name, user.Email, user.Password.Hash, user.Activated).Scan(&user.ID, &user.Created_at); err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return internal.ErrDuplicateEmail
		}
		return err
	}

	return nil
}
