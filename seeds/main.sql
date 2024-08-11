INSERT INTO days (
	date, 
	lunch, 
	x_period, 
	rotation_day, 
	location, 
	notes, 
	ap_info, 
	cc_info, 
	grade_9, 
	grade_10, 
	grade_11, 
	grade_12
) VALUES 
(
	'2024-12-18', 
	'lunch', 
	'x_period', 
	'rotation_day', 
	'location',
	'notes', 
	'ap_info', 
	'cc_info', 
	'grade_9', 
	'grade_10', 
	'grade_11', 
	'grade_12'
),
(
	'2024-12-19', 
	'lunch', 
	'x_period', 
	'rotation_day', 
	'location',
	'notes', 
	'ap_info', 
	'cc_info', 
	'grade_9', 
	'grade_10', 
	'grade_11', 
	'grade_12'
),
(
	'2024-12-20', 
	'lunch', 
	'x_period', 
	'rotation_day', 
	'location',
	'notes', 
	'ap_info', 
	'cc_info', 
	'grade_9', 
	'grade_10', 
	'grade_11', 
	'grade_12'
);

INSERT INTO users (
	id, 
	name, 
	email,
	password_hash, 
	is_admin,
	status
) VALUES (
	'aea38951-ca26-4e76-ad65-d5296a0095e6',
	'name',
	'email',
	'password_hash',
	false,
	'unverified'	
);

INSERT INTO email_verifications (
	id, 
	user_id
) VALUES 
	('df27db6d-14b4-4fe9-a975-e700e9135733', 'aea38951-ca26-4e76-ad65-d5296a0095e6'),
	('01b7050f-40fe-4d25-bea4-b72b29134e75', 'aea38951-ca26-4e76-ad65-d5296a0095e6'),
	('3505a964-26da-465d-8b26-1c371dbe04c5', 'aea38951-ca26-4e76-ad65-d5296a0095e6');
