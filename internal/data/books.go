package data

import (
	"context"
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

func (b BookClub) InsertBook(book *Book) error {
	// err := b.DoesProductExists(id)

	// if err != nil {
	// 	return err
	// }

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
	
	INSERT INTO books (title, isbn, publication_date, genre, description) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id;
	
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

// func (p ProductModel) GetReview(id int64, rid int64) (*Reviews, error) {
// 	if id < 1 {
// 		return nil, ErrRecordNotFound
// 	} else if rid < 1 {
// 		return nil, ErrRecordNotFound
// 	}

// 	query := `
// 	SELECT R.id, P.name, R.rating, R.helpful_count, R.comment, R.created_at, R.updated_at
// 	FROM reviews AS R
// 	INNER JOIN products AS P ON P.id = R.product_id
// 	WHERE R.id = $1 AND R.product_id = $2;

// 	`

// 	args := []any{rid, id}

// 	var review Reviews

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&review.ID, &review.Product_Id, &review.Rating, &review.Helpful_Count, &review.Comment, &review.Created_at, &review.Updated_at)

// 	if err != nil {
// 		switch {
// 		case errors.Is(err, sql.ErrNoRows):
// 			return nil, ErrRecordNotFound
// 		default:
// 			return nil, err
// 		}
// 	}
// 	return &review, nil

// }

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

// func (p ProductModel) GetAllReviews(product int64, filters Filters) ([]*Reviews, error) {
// 	query := fmt.Sprintf(`
// 	SELECT R.id, P.name, R.rating, R.helpful_count, R.comment, R.created_at, R.updated_at
// 	FROM reviews AS R
// 	INNER JOIN products AS P ON P.id = R.product_id
// 	WHERE P.id = $1 OR NOT EXISTS (SELECT 1 FROM products WHERE id = $1)
// 	ORDER BY %s %s, R.id ASC
// 	LIMIT $2 OFFSET $3
// 	`, filters.sortColumn(), filters.sortDirection())

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	rows, err := p.DB.QueryContext(ctx, query, product, filters.limit(), filters.offset())
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	reviews := []*Reviews{}

// 	for rows.Next() {
// 		var review Reviews
// 		err := rows.Scan(&review.ID, &review.Product_Id, &review.Rating, &review.Helpful_Count, &review.Comment, &review.Created_at, &review.Updated_at)

// 		if err != nil {
// 			return nil, err
// 		}

// 		reviews = append(reviews, &review)
// 	}

// 	err = rows.Err()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return reviews, nil
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
