package books

import "github.com/lib/pq"

func (b *BookModel) Get(id int64) (*Book, error) {
	query := `
	SELECT
		id,
		title,
		 
		author,
		release_date,
		publisher,
		country ,
		isbn ,
		pages ,
		language,
		description,
		genres,
		cover_image_url,
		price
	FROM 
		books 
	WHERE
		id = $1`
	book := &Book{}
	err := b.DB.QueryRow(query, id).Scan(
		&book.ID, &book.Title,
		&book.Author, &book.Release_date, &book.Publisher,
		&book.Country, &book.Isbn, &book.Pages, &book.Language,
		&book.Description, pq.Array(&book.Genres), &book.CoverImageUrl, &book.Price,
	)
	return book, err
}

func (b *BookModel) Latest() ([]*Book, error) {
	query := `
	SELECT 
		id,
		title,
		 
		author,
		release_date,
		publisher,
		country ,
		isbn ,
		pages ,
		language,
		description,
		genres,
		cover_image_url,
		price
	FROM 
		books 
	ORDER BY 
		id ASC
	FROM 
		books 
	ORDER BY 
		id DESC LIMIT 10`

	books := []*Book{}

	rows, err := b.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		book := &Book{}

		err := rows.Scan(
			&book.ID, &book.Title,
			&book.Author, &book.Release_date, &book.Publisher,
			&book.Country, &book.Isbn, &book.Pages, &book.Language,
			&book.Description, pq.Array(&book.Genres), &book.CoverImageUrl, &book.Price,
		)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (b *BookModel) GetAll() (*[]Book, error) {
	query := `
	SELECT 
		id,
		title,
		
		author,
		release_date,
		publisher,
		country,
		isbn,
		pages,
		language,
		description,
		genres,
		cover_image_url,
		price
	FROM 
		books 
	ORDER BY 
		id ASC
		`

	books := []Book{}

	rows, err := b.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var book Book

		err := rows.Scan(
			&book.ID, &book.Title, &book.Author,
			&book.Release_date, &book.Publisher,
			&book.Country, &book.Isbn, &book.Pages, &book.Language,
			&book.Description, pq.Array(&book.Genres), &book.CoverImageUrl, &book.Price,
		)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return &books, nil
}
