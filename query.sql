-- name: CreateBook :one
INSERT INTO books (
  name, link, image, status
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteBook :exec
DELETE FROM books
WHERE id = $1;

-- name: UpdateBook :exec
UPDATE books
  set name = $2,
  link = $3,
  image = $4,
  status = $5
WHERE id = $1;

-- name: GetReadingBook :many
SELECT name, link, image
FROM books
WHERE status = 1;

-- name: GetDoneBooks :many
SELECT name, link, image
FROM books
WHERE status = 2;

-- name: GetWantToReadBooks :many
SELECT name, link, image FROM books
WHERE status = 3;