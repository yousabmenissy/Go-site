package tokens

func (t *TokenModel) Insert(token *Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope) VALUES ($1, $2, $3, $4)`

	_, err := t.DB.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	return err
}
