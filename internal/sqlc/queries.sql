-- name: CreateUser :one
INSERT INTO users (username, email, password)
VALUES ($1, $2, $3)
RETURNING id, username, email, password;

-- name: GetUserByUsername :one
SELECT id, username, email, password
FROM users
WHERE username = $1;

-- name: GetUserByID :one
SELECT id, username, email, password
FROM users
WHERE id = $1;

-- name: CreateDocument :one
INSERT INTO documents (user_id, content, private)
VALUES ($1, $2, $3)
RETURNING id, user_id, content, private, created, updated;

-- name: GetDocuments :many
SELECT id, user_id, content, private, created, updated
FROM documents
ORDER BY updated DESC
LIMIT $1;

-- name: GetDocumentsByUserID :many
SELECT id, user_id, content, private, created, updated
FROM documents
WHERE user_id = $1
ORDER BY updated DESC;

-- name: UpdateDocument :one
UPDATE documents
SET content = $2, updated = now()
WHERE id = $1
RETURNING id, user_id, content, private, created, updated;

-- name: DeleteDocument :exec
DELETE FROM documents
WHERE id = $1;