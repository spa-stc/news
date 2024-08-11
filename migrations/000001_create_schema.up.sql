CREATE TABLE IF NOT EXISTS days (
	date DATE PRIMARY KEY,

	lunch VARCHAR NOT NULL, 
	x_period VARCHAR NOT NULL,
	rotation_day VARCHAR NOT NULL,
	location VARCHAR NOT NULL, 
	notes VARCHAR NOT NULL, 
	ap_info VARCHAR NOT NULL,
	cc_info VARCHAR NOT NULL, 
	grade_9 VARCHAR NOT NULL, 
	grade_10 VARCHAR NOT NULL, 
	grade_11 VARCHAR NOT NULL, 
	grade_12 VARCHAR NOT NULL, 
	
	created_ts TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_ts TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TYPE user_status AS ENUM (
	'verified', 
	'unverified', 
	'banned'
);

CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY NOT NULL,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL UNIQUE,

	password_hash VARCHAR NOT NULL, 

	is_admin BOOLEAN NOT NULL DEFAULT false,
	status user_status NOT NULL DEFAULT 'unverified',

	created_ts TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_ts TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS email_verifications (
	id UUID PRIMARY KEY NOT NULL,

	used BOOLEAN NOT NULL DEFAULT FALSE,
	user_id UUID NOT NULL REFERENCES users (id), 

	created_ts TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX email_verifications_users_idx ON email_verifications (user_id);
