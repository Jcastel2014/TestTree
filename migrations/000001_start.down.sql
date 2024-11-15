SELECT B.id, B.title, B.isbn, A.name AS author, B.publication_date, B.genre, B.description, B.average_rating
FROM books AS B
INNER JOIN book_authors AS BA 
ON B.id = BA.book_id
INNER JOIN authors AS A 
ON A.id = BA.author_id

	WHERE P.id = $1 OR NOT EXISTS (SELECT 1 FROM products WHERE id = $1)
