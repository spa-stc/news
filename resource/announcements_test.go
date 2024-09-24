package resource_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"stpaulacademy.tech/newsletter/resource"
	"stpaulacademy.tech/newsletter/util/sliceutil"
	"stpaulacademy.tech/newsletter/util/testutil"
)

func TestAnnouncementsResource(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("test_range_not_found", func(t *testing.T) {
		tx := testutil.TestTx(ctx, t)
		_, err := resource.GetManyAnnouncementsByCurrentDay(
			ctx,
			tx,
			testutil.ParseDate(t, "2023-01-01"),
		)

		require.NoError(t, err)
	})

	t.Run("test_range", func(t *testing.T) {
		tx := testutil.TestTx(ctx, t)

		expected := []resource.Announcement{
			{
				ID:           1,
				Title:        "title",
				Author:       "author",
				Content:      "content",
				CreatedTS:    time.UnixMicro(0),
				UpdatedTS:    time.UnixMicro(0),
				DisplayStart: time.UnixMicro(0),
				DisplayEnd:   time.UnixMicro(0),
			},
		}

		a, err := resource.GetManyAnnouncementsByCurrentDay(
			ctx,
			tx,
			testutil.ParseDate(t, "2024-12-18"),
		)
		require.NoError(t, err)

		a = sliceutil.Map(a, func(e resource.Announcement) resource.Announcement {
			e.CreatedTS = time.UnixMicro(0)
			e.UpdatedTS = time.UnixMicro(0)
			e.DisplayEnd = time.UnixMicro(0)
			e.DisplayStart = time.UnixMicro(0)

			return e
		})

		require.Equal(t, expected, a)
	})

	t.Run("test_insert", func(t *testing.T) {
		tx := testutil.TestTx(ctx, t)

		expected := resource.Announcement{
			ID:           3,
			Title:        "title",
			Author:       "author",
			Content:      "content",
			CreatedTS:    time.UnixMicro(0).UTC(),
			UpdatedTS:    time.UnixMicro(0).UTC(),
			DisplayStart: time.UnixMicro(0).UTC(),
			DisplayEnd:   time.UnixMicro(0).UTC(),
		}

		a, err := resource.InsertAnnouncement(ctx, tx, resource.NewAnnouncement{
			Title:        expected.Title,
			Author:       expected.Author,
			Content:      expected.Content,
			DisplayStart: expected.DisplayStart,
			DisplayEnd:   expected.DisplayEnd,
		})
		require.NoError(t, err)

		a.CreatedTS = time.UnixMicro(0).UTC()
		a.UpdatedTS = time.UnixMicro(0).UTC()

		require.Equal(t, expected, a)
	})

	t.Run("test_delete", func(t *testing.T) {
		tx := testutil.TestTx(ctx, t)

		a, err := resource.InsertAnnouncement(ctx, tx, resource.NewAnnouncement{
			Title:        "Test",
			Author:       "E",
			Content:      "Hi",
			DisplayStart: time.UnixMicro(0),
			DisplayEnd:   time.UnixMicro(0),
		})
		require.NoError(t, err)

		err = resource.DeleteAnnouncement(ctx, tx, a.ID)

		require.NoError(t, err)
	})

	t.Run("test_get_upcoming", func(t *testing.T) {
		tx := testutil.TestTx(ctx, t)

		expected := []resource.Announcement{
			{
				ID:           1,
				Title:        "title",
				Author:       "author",
				Content:      "content",
				CreatedTS:    time.UnixMicro(0),
				UpdatedTS:    time.UnixMicro(0),
				DisplayStart: time.UnixMicro(0),
				DisplayEnd:   time.UnixMicro(0),
			},
			{
				ID:           2,
				Title:        "title",
				Author:       "author",
				Content:      "content",
				CreatedTS:    time.UnixMicro(0),
				UpdatedTS:    time.UnixMicro(0),
				DisplayStart: time.UnixMicro(0),
				DisplayEnd:   time.UnixMicro(0),
			},
		}

		a, err := resource.GetUpcomingAnnouncements(ctx, tx, testutil.ParseDate(t, "2023-01-01"))
		require.NoError(t, err)

		a = sliceutil.Map(a, func(e resource.Announcement) resource.Announcement {
			e.CreatedTS = time.UnixMicro(0)
			e.UpdatedTS = time.UnixMicro(0)
			e.DisplayEnd = time.UnixMicro(0)
			e.DisplayStart = time.UnixMicro(0)

			return e
		})

		require.Equal(t, expected, a)
	})
}
