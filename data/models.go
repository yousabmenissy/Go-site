package data

import (
	"database/sql"
	"site/data/books"
	"site/data/tokens"
	"site/data/users"
)

type Models struct {
	Books  *books.BookModel
	Users  *users.UserModel
	Tokens *tokens.TokenModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Books:  &books.BookModel{DB: db},
		Users:  &users.UserModel{DB: db},
		Tokens: &tokens.TokenModel{DB: db},
	}
}
