package books

import (
	"database/sql"
	"time"
)

type BookModel struct {
	DB *sql.DB
}

type Book struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Release_date  time.Time `json:"release_date"`
	Publisher     string    `json:"publisher"`
	Country       string    `json:"country"`
	Isbn          string    `json:"isbn"`
	Pages         int       `json:"pages"`
	Language      string    `json:"language"`
	Description   string    `json:"description"`
	Genres        []string  `json:"genres"`
	CoverImageUrl string    `json:"cover_image_url"`
	Price         float64   `json:"price"`
}
