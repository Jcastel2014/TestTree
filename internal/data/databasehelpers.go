package data

import (
	"context"
	"time"

	"github.com/Jcastel2014/test3/internal/validator"
)

func (b BookClub) DoesAuthorExists(author string) (error, int) {
	query := `
		SELECT id
		FROM authors
		WHERE name = $1
	`
	args := []any{author}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int

	err := b.DB.QueryRowContext(ctx, query, args...).Scan(&id)

	if err != nil {
		return err, -1
	}

	if id < 1 {
		return nil, 0
	}

	return nil, id
}

func ValidateBook(v *validator.Validator, book *Book) {

	v.Check(book.Title != "", "title", "must be provided")
	v.Check(len(book.Title) <= 255, "title", "must not be more than 100 byte long")

	v.Check(book.ISBN != "", "isbn", "must be provided")
	v.Check(len(book.ISBN) == 13, "isbn", "must be exactly 13 characters long")

	v.Check(book.Author != "", "author", "must be provided")
	v.Check(len(book.Author) <= 100, "author", "must not be more than 100 characters long")

	v.Check(book.Genre != "", "genre", "must be provided")
	v.Check(len(book.Genre) <= 50, "genre", "must not be more than 50 characters long")

	v.Check(len(book.Description) <= 1000, "description", "must not be more than 1000 characters long")

	v.Check(!book.Publication_Date.IsZero(), "publication_date", "must be provided")
	v.Check(book.Publication_Date.Before(time.Now()), "publication_date", "must not be in the future")

	// v.Check(review.Rating > 0, "rating", "must be greater than 0")
	// v.Check(review.Rating <= 5, "rating", "must be less than 5")
	// v.Check(len(review.Comment) <= 100, "comment", "must not be more than 100 byte long")

}
