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
SELECT 
	users.*, 
	COUNT(email_verifications.id) AS verification_attempts 
FROM 
	users LEFT JOIN email_verifications ON user_id = users.id 
WHERE 
	users.id = @id
GROUP BY 
	users.id
LIMIT 1; 

-- name: GetUserByEmail :one
SELECT 
	users.*, 
	COUNT(email_verifications.id) AS verification_attempts 
FROM 
	users LEFT JOIN email_verifications ON user_id = users.id 
WHERE 
	users.email = @email
GROUP BY 
	users.id 
LIMIT 1; 

-- name: UpdateUserByID :exec
UPDATE 
	users 
SET 
	status = CASE WHEN @status_do_update::boolean 
		THEN @status::user_status ELSE status END,

	password_hash = CASE WHEN @password_hash_do_update::boolean
		THEN @password_hash::VARCHAR ELSE password_hash END, 

	updated_ts = NOW()
WHERE 
	id = @id;

-- name: GetTokenClaimsByUserID :one
SELECT id, is_admin, status FROM users WHERE id = @id LIMIT 1;
