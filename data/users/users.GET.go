package users

import (
	"database/sql"
	"errors"
	"site/internal"
)

func (u *UserModel) GetUserByID(id int64) (*User, error) {
	query := `SELECT 
		id,
		name, 
		email, 
		password_hash, 
		activated, 
		created_at
	FROM 
		users 
	WHERE 
		id = $1`

	var user User
	err := u.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Created_at,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, internal.ErrNoRecord

		default:
			return nil, err
		}
	}

	return &user, nil

}

func (u *UserModel) GetUserByEmail(email string) (*User, error) {
	query := `SELECT 
	id,
	name, 
	email, 
	password_hash, 
	activated, 
	created_at
FROM 
	users 
WHERE 
	email = $1`

	var user User
	err := u.DB.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password.Hash, &user.Activated, &user.Created_at)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, internal.ErrNoRecord

		default:
			return nil, err
		}
	}

	return &user, nil
}
