-- name: GetEmailVerificationBySecret :one 
SELECT secret, user_id, used FROM email_verifications WHERE secret = @secret LIMIT 1;

-- name: InsertEmailVerification :one 
INSERT INTO email_verifications (
	id,
	secret,
	user_id
) VALUES (
	@id,
	@secret, 
	@user_id 
) RETURNING secret, user_id, used;
