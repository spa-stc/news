-- Users Table:
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY,
	email TEXT UNIQUE NOT NULL,
	name TEXT NOT NULL,
	password TEXT NOT NULL,
	is_admin BOOLEAN NOT NULL DEFAULT FALSE,
	email_verified BOOLEAN NOT NULL DEFAULT FALSE, 
	disabled BOOLEAN NOT NULL DEFAULT FALSE,
	created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
	updated_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now'))
);

CREATE INDEX users_email_idx ON users (email);
