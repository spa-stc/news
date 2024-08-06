-- name: GetManyDaysByDate :many
SELECT
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
	grade_12,
	created_ts,
	updated_ts
FROM 
	days 
WHERE
	date = ANY(@dates::date[]);

-- name: UpsertDay :one
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
) VALUES (
	@date, 
	@lunch, 
	@x_period, 
	@rotation_day, 
	@location, 
	@notes, 
	@ap_info,
	@cc_info, 
	@grade_9,
	@grade_10,
	@grade_11,
	@grade_12
) 
ON CONFLICT DO UPDATE 
SET 
	lunch = @lunch,
	x_period = @x_period,
	rotation_day = @rotation_day,
	location = @location,
	notes = @notes,
	ap_info = @ap_info,
	cc_info = @cc_info, 
	grade_9 = @grade_9,
	grade_10 = @grade_10,
	grade_11 = @grade_11,
	grade_12 = @grade_12
RETURNING *;
