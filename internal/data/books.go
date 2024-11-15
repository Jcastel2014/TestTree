package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Book struct {
	ID               int64     `json:"id"`
	Title            string    `json:"title"`
	ISBN             string    `json:"isbn"`
	Author           string    `json:"author"`
	Genre            string    `json:"genre"`
	Description      string    `json:"description"`
	Publication_Date time.Time `json:"created_at"`
	Average_rating   float64   `json:"average_rating"`
}

func (b BookClub) GetAllBooks(filters Filters) ([]*Book, error) {
	query := fmt.Sprintf(`
	SELECT B.id, B.title, B.isbn, A.name AS author, B.publication_date, B.genre, B.description, B.average_rating
	FROM books AS B
	INNER JOIN book_authors AS BA 
	ON B.id = BA.book_id
	INNER JOIN authors AS A 
	ON A.id = BA.author_id
	ORDER BY %s %s, B.id ASC
	LIMIT $1 OFFSET $2
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := b.DB.QueryContext(ctx, query, filters.limit(), filters.offset())
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	books := []*Book{}

	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.ISBN, &book.Author, &book.Publication_Date, &book.Genre, &book.Description, &book.Average_rating)
		if err != nil {
			return nil, err
		}

		books = append(books, &book)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (b BookClub) InsertBook(book *Book) error {

	var idA int

	query := `
	INSERT INTO authors(name)
	VALUES ($1)
	RETURNING id
	
	`

	err, idA := b.DoesAuthorExists(book.Author)

	if err != nil {
		return err
	}

	if idA == 0 {

		args := []any{book.Author}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err = b.DB.QueryRowContext(ctx, query, args...).Scan(&idA)

	}

	// args := []any{book.Author}
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// err = b.DB.QueryRowContext(ctx, query, args...).Scan(&idA)

	if err != nil {
		return err
	}

	query = `
	
	INSERT INTO books (title, isbn, publication_date, genre, description, average_rating) 
	VALUES ($1, $2, $3, $4, $5, 0) RETURNING id;
	
	`

	var idB int
	args := []any{book.Title, book.ISBN, book.Publication_Date, book.Genre, book.Description}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = b.DB.QueryRowContext(ctx, query, args...).Scan(&idB)

	if err != nil {
		return err
	}

	query = `
		INSERT INTO book_authors (book_id, author_id) 
		VALUES ($1, $2) RETURNING id;

	`

	args = []any{idB, idA}
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = b.DB.QueryRowContext(ctx, query, args...).Scan(&idB)

	if err != nil {
		return err
	}

	return err

	// return b.UpdateAverage(id)

}

// func (p ProductModel) UpdateAverage(pid int64) error {

// 	query := `
// 	UPDATE products
// 	Set average_rating = (select AVG(rating)::NUMERIC(10,2) from reviews where product_id = $1)
// 	WHERE id = $1
// 	RETURNING id
// 	`

// 	args := []any{pid}

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return p.DB.QueryRowContext(ctx, query, args...).Scan(&pid)
// }

func (b BookClub) GetBook(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
	SELECT B.id, B.title, B.isbn, A.name AS author, B.publication_date, B.genre, B.description, B.average_rating
	FROM books AS B
	INNER JOIN book_authors AS BA 
	ON B.id = BA.book_id
	INNER JOIN authors AS A 
	ON A.id = BA.author_id
	WHERE B.id = $1

	`

	args := []any{id}

	var book Book

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID, &book.Title, &book.ID, &book.Author, &book.Publication_Date, &book.Genre, &book.Description, &book.Average_rating)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &book, nil

}

// func (p ProductModel) UpdateReview(review *Reviews, id int64) error {
// 	query := `
// 	UPDATE reviews
// 	SET rating =$1, comment=$2, updated_at=$3
// 	WHERE id = $4
// 	RETURNING product_id
// 	`

// 	args := []any{review.Rating, review.Comment, time.Now(), id}
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&review.ID)

// 	if err != nil {
// 		return err
// 	}

// 	return p.UpdateAverage(review.ID)

// }

// func (p ProductModel) DeleteReview(id int64, rid int64) error {
// 	err := p.DoesProductExists(id)

// 	if err != nil {
// 		return err
// 	}

// 	query := `
// 	DELETE FROM reviews
// 	WHERE ID = $1
// 	`

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	result, err := p.DB.ExecContext(ctx, query, rid)
// 	if err != nil {
// 		return err
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}

// 	if rowsAffected == 0 {
// 		return ErrRecordNotFound
// 	}

// 	return p.UpdateAverage(id)

// }

// func (p ProductModel) DoesProductExists(id int64) error {
// 	query := `
// 		SELECT COUNT(*)
// 		FROM products
// 		WHERE id = $1
// 	`
// 	args := []any{id}

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	var count int

// 	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&count)

// 	if err != nil {
// 		return err
// 	}

// 	if count < 1 {
// 		return errors.New("no rows found")
// 	}

// 	return nil
// }
