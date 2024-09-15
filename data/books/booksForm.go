package books

import (
	"site/internal/validation"
	"strconv"
	"strings"
	"time"
)

type BookSubmitForm struct {
	Input struct {
		Title         string `json:"title"`
		Author        string `json:"author"`
		Release_date  string `json:"release_date"`
		Publisher     string `json:"publisher"`
		Country       string `json:"country"`
		Isbn          string `json:"isbn"`
		Pages         string `json:"pages"`
		Language      string `json:"language"`
		Description   string `json:"description"`
		Genres        string `json:"genres"`
		CoverImageUrl string `json:"cover_image_url"`
		Price         string `json:"price"`
	}
	Output Book
	V      validation.Validation
}

func (b *BookSubmitForm) Prepare() error {
	b.Validate()

	b.Output.Title = b.Input.Title
	b.Output.Author = b.Input.Author
	newdate, err := time.Parse("2006-1-2", b.Input.Release_date)
	if err != nil {
		return err
	}
	if newdate.After(time.Now()) {
		b.V.AddError("release_date", "invalid date")
	}
	b.Output.Release_date = newdate
	b.Output.Publisher = b.Input.Publisher
	b.Output.Country = b.Input.Country
	b.Output.Isbn = b.Input.Isbn
	pages, err := strconv.Atoi(b.Input.Pages)
	if err != nil {
		return err
	}
	if pages < 0 {
		b.V.AddError("pages", "invalid pages field")
	}
	b.Output.Pages = pages
	b.Output.Language = b.Input.Language
	b.Output.Description = b.Input.Description
	b.Output.Genres = strings.Split(strings.TrimSpace(b.Input.Genres), ",")
	b.Output.CoverImageUrl = b.Input.CoverImageUrl
	price, err := strconv.ParseFloat(strings.TrimSpace(b.Input.Price), 64)
	if err != nil {
		return err
	}
	if price < 0 {
		b.V.AddError("price", "invalid price field")
	}
	pricestr := strconv.FormatFloat(price, 'f', 2, 64)
	b.Output.Price, _ = strconv.ParseFloat(pricestr, 64)

	return nil
}

func (b *BookSubmitForm) Validate() {
	b.V.Check(validation.NotBlank(b.Input.Title), "title", "required")
	b.V.Check(validation.NotBlank(b.Input.Author), "author", "required")
	b.V.Check(validation.NotBlank(b.Input.Release_date), "release_date", "required")
	b.V.Check(validation.NotBlank(b.Input.Publisher), "publisher", "required")
	b.V.Check(validation.NotBlank(b.Input.Country), "country", "required")
	b.V.Check(validation.NotBlank(b.Input.Isbn), "isbn", "required")
	b.V.Check(validation.NotBlank(b.Input.Pages), "pages", "required")
	b.V.Check(validation.NotBlank(b.Input.Language), "language", "required")
	b.V.Check(validation.NotBlank(b.Input.Description), "description", "required")
	b.V.Check(validation.NotBlank(b.Input.Genres), "genres", "required")
	b.V.Check(validation.NotBlank(b.Input.CoverImageUrl), "cover_image_url", "required")
	b.V.Check(validation.NotBlank(b.Input.Price), "price", "required")
}
