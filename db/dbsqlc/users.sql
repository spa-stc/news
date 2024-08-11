-- name: InsertUser :one
INSERT INTO users (
	id,
	name,
	email,
	password_hash
) VALUES (
	@id,
	@name,
	@email,
	@password_hash
) RETURNING *;

-- name: GetUserByID :one 
SELECT * FROM users WHERE id = @id;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = @email;

