package resource

import (
	"context"
	"time"

	"stpaulacademy.tech/newsletter/db"
	"stpaulacademy.tech/newsletter/db/dbsqlc"
	"stpaulacademy.tech/newsletter/util/sliceutil"
)

type Announcement struct {
	ID int64

	Title        string
	Author       string
	Content      string
	DisplayStart time.Time
	DisplayEnd   time.Time

	CreatedTS time.Time
	UpdatedTS time.Time
}

func fromSqlcAnnouncement(d dbsqlc.Announcement) Announcement {
	return Announcement{
		ID:           d.ID,
		Title:        d.Title,
		Author:       d.Author,
		Content:      d.Content,
		DisplayStart: d.DisplayStart.UTC(),
		DisplayEnd:   d.DisplayEnd.UTC(),

		CreatedTS: d.CreatedTs.UTC(),
		UpdatedTS: d.UpdatedTs.UTC(),
	}
}

func GetManyAnnouncementsByCurrentDay(ctx context.Context, e db.Executor, today time.Time) ([]Announcement, error) {
	sqlc := dbsqlc.New()

	announcements, err := sqlc.GetManyAnnouncementsWhereInRange(ctx, e, today)
	if err != nil {
		return nil, db.HandleError(err)
	}

	return sliceutil.Map(announcements, fromSqlcAnnouncement), nil
}
