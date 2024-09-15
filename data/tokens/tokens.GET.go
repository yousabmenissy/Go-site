package tokens

import (
	"database/sql"
	"errors"
	"site/internal"
)

func (t *TokenModel) GetUserToken(id int64) ([]byte, error) {
	query := `SELECT hash FROM tokens WHERE user_id = $1 AND expiry > NOW()`

	var tokenHash []byte
	err := t.DB.QueryRow(query, id).Scan(&tokenHash)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, internal.ErrNoRecord

		default:
			return nil, err
		}
	}

	return tokenHash, nil
}
