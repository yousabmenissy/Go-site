package books

import "github.com/lib/pq"

func (b *BookModel) Insert(book *Book) error {
	query := `
	INSERT INTO 
		books (
		"title",
		"author",
		"release_date",
		"publisher",
		"country",
		"isbn",
		"pages",
		"language",
		"description",
		"genres",
		"cover_image_url",
		"price"
		) 
		VALUES (
		$1, $2, $3, 
		$4, $5, $6, 
		$7, $8, $9, 
		$10, $11, $12
		) 
		 RETURNING id`
	err := b.DB.QueryRow(
		query,
		book.Title,
		book.Author,
		book.Release_date,
		book.Publisher,
		book.Country,
		book.Isbn,
		book.Pages,
		book.Language,
		book.Description,
		pq.Array(book.Genres),
		book.CoverImageUrl,
		book.Price,
	).Scan(&book.ID)
	return err
}
