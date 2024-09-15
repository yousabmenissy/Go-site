package tokens

func (t *TokenModel) DeleteAllForUser(id int64, scope string) error {
	query := `DELETE FROM tokens WHERE user_id = $1 AND scope = $2`

	_, err := t.DB.Exec(query, id, scope)
	return err
}
