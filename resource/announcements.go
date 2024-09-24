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

type NewAnnouncement struct {
	Title        string
	Author       string
	Content      string
	DisplayStart time.Time
	DisplayEnd   time.Time
}

func toSqlcNewAnnouncement(n NewAnnouncement) dbsqlc.InsertAnnouncementParams {
	return dbsqlc.InsertAnnouncementParams{
		Title:        n.Title,
		Author:       n.Author,
		Content:      n.Content,
		DisplayStart: n.DisplayStart.UTC(),
		DisplayEnd:   n.DisplayEnd.UTC(),
	}
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

func InsertAnnouncement(ctx context.Context, e db.Executor, newAnnouncement NewAnnouncement) (Announcement, error) {
	announcement, err := dbsqlc.New().InsertAnnouncement(ctx, e, toSqlcNewAnnouncement(newAnnouncement))
	if err != nil {
		return Announcement{}, db.HandleError(err)
	}

	return fromSqlcAnnouncement(announcement), nil
}

func DeleteAnnouncement(ctx context.Context, e db.Executor, id int64) error {
	err := dbsqlc.New().DeleteAnnouncement(ctx, e, id)
	if err != nil {
		return db.HandleError(err)
	}

	return nil
}

func GetUpcomingAnnouncements(ctx context.Context, e db.Executor, today time.Time) ([]Announcement, error) {
	announcements, err := dbsqlc.New().GetAllUpcomingAnnouncements(ctx, e, today)
	if err != nil {
		return nil, db.HandleError(err)
	}

	return sliceutil.Map(announcements, fromSqlcAnnouncement), nil
}
