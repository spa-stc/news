-- name: GetManyAnnouncementsWhereInRange :many 
SELECT * FROM announcements WHERE display_start <= @today AND display_end >= @today;
