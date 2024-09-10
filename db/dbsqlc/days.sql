-- name: GetManyDaysByDate :many
SELECT
	date, 
	lunch,
	x_period, 
	rotation_day,
	location, 
	notes,
	cc_info,
	grade_9,
	grade_10,
	grade_11, 
	grade_12,
	created_ts,
	updated_ts
FROM 
	days 
WHERE
	date = ANY(@dates::date[]);

-- name: GetManyDaysbyDateRange :many 
SELECT * FROM days WHERE date >= $1 AND date <= $2;

-- name: UpsertDay :one
INSERT INTO days (
	date, 
	lunch,
	x_period, 
	rotation_day,
	location, 
	notes,
	cc_info,
	grade_9,
	grade_10,
	grade_11, 
	grade_12
) VALUES (
	@date, 
	@lunch, 
	@x_period, 
	@rotation_day, 
	@location, 
	@notes, 
	@cc_info, 
	@grade_9,
	@grade_10,
	@grade_11,
	@grade_12
) 
ON CONFLICT (date) DO UPDATE 
SET 
	lunch = @lunch,
	x_period = @x_period,
	rotation_day = @rotation_day,
	location = @location,
	notes = @notes,
	cc_info = @cc_info, 
	grade_9 = @grade_9,
	grade_10 = @grade_10,
	grade_11 = @grade_11,
	grade_12 = @grade_12
RETURNING *;
