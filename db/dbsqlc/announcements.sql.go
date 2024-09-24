// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: announcements.sql

package dbsqlc

import (
	"context"
	"time"
)

const deleteAnnouncement = `-- name: DeleteAnnouncement :exec
DELETE FROM announcements WHERE id = $1
`

func (q *Queries) DeleteAnnouncement(ctx context.Context, db DBTX, id int64) error {
	_, err := db.Exec(ctx, deleteAnnouncement, id)
	return err
}

const getAllUpcomingAnnouncements = `-- name: GetAllUpcomingAnnouncements :many
SELECT id, title, author, content, display_start, display_end, created_ts, updated_ts FROM announcements WHERE display_end >= $1
`

func (q *Queries) GetAllUpcomingAnnouncements(ctx context.Context, db DBTX, today time.Time) ([]Announcement, error) {
	rows, err := db.Query(ctx, getAllUpcomingAnnouncements, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Announcement
	for rows.Next() {
		var i Announcement
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Author,
			&i.Content,
			&i.DisplayStart,
			&i.DisplayEnd,
			&i.CreatedTs,
			&i.UpdatedTs,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getManyAnnouncementsWhereInRange = `-- name: GetManyAnnouncementsWhereInRange :many
SELECT id, title, author, content, display_start, display_end, created_ts, updated_ts FROM announcements WHERE display_start <= $1 AND display_end >= $1
`

func (q *Queries) GetManyAnnouncementsWhereInRange(ctx context.Context, db DBTX, today time.Time) ([]Announcement, error) {
	rows, err := db.Query(ctx, getManyAnnouncementsWhereInRange, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Announcement
	for rows.Next() {
		var i Announcement
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Author,
			&i.Content,
			&i.DisplayStart,
			&i.DisplayEnd,
			&i.CreatedTs,
			&i.UpdatedTs,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertAnnouncement = `-- name: InsertAnnouncement :one
INSERT INTO announcements (
	title,
	author,
	content, 
	display_start, 
	display_end
) VALUES (
	$1, 
	$2,
	$3, 
	$4, 
	$5
) RETURNING id, title, author, content, display_start, display_end, created_ts, updated_ts
`

type InsertAnnouncementParams struct {
	Title        string
	Author       string
	Content      string
	DisplayStart time.Time
	DisplayEnd   time.Time
}

func (q *Queries) InsertAnnouncement(ctx context.Context, db DBTX, arg InsertAnnouncementParams) (Announcement, error) {
	row := db.QueryRow(ctx, insertAnnouncement,
		arg.Title,
		arg.Author,
		arg.Content,
		arg.DisplayStart,
		arg.DisplayEnd,
	)
	var i Announcement
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Author,
		&i.Content,
		&i.DisplayStart,
		&i.DisplayEnd,
		&i.CreatedTs,
		&i.UpdatedTs,
	)
	return i, err
}
