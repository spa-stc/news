-- name: GetManyAnnouncementsWhereInRange :many 
SELECT * FROM announcements WHERE display_start <= @today AND display_end >= @today;

-- name: GetAllUpcomingAnnouncements :many 
SELECT * FROM announcements WHERE display_end >= @today;

-- name: InsertAnnouncement :one
INSERT INTO announcements (
	title,
	author,
	content, 
	display_start, 
	display_end
) VALUES (
	@title, 
	@author,
	@content, 
	@display_start, 
	@display_end
) RETURNING *;

-- name: DeleteAnnouncement :exec
DELETE FROM announcements WHERE id = @id;

