package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"stpaulacademy.tech/newsletter/api"
	"stpaulacademy.tech/newsletter/util/sliceutil"
	"stpaulacademy.tech/newsletter/util/testutil"
)

func TestGetWeek(t *testing.T) {
	t.Parallel()

	ctx := testutil.Setup(t)
	srv := NewTestAPI(ctx, t, time.Date(2024, 12, 18, 0, 0, 0, 0, time.UTC))

	t.Run("test_success", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/week", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		srv.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var res api.WeekResponse
		err = json.NewDecoder(rr.Body).Decode(&res)
		require.NoError(t, err)

		dates := []time.Time{
			testutil.ParseDate(t, "2024-12-18"),
			testutil.ParseDate(t, "2024-12-19"),
			testutil.ParseDate(t, "2024-12-20"),
		}

		expected := sliceutil.Map(dates, func(s time.Time) api.Day {
			return api.Day{
				Date:        s,
				Lunch:       "lunch",
				XPeriod:     "x_period",
				RotationDay: "rotation_day",
				Location:    "location",
				Notes:       "notes",
				ApInfo:      "ap_info",
				CcInfo:      "cc_info",
				Grade9:      "grade_9",
				Grade10:     "grade_10",
				Grade11:     "grade_11",
				Grade12:     "grade_12",
				CreatedTS:   time.UnixMicro(0),
				UpdatedTS:   time.UnixMicro(0),
			}
		})

		found := sliceutil.Map(res.Days, func(t api.Day) api.Day {
			t.CreatedTS = time.UnixMicro(0)
			t.UpdatedTS = time.UnixMicro(0)

			return t
		})

		require.ElementsMatch(t, expected, found)
	})

	t.Run("test_failure", func(t *testing.T) {
		srv := NewTestAPI(ctx, t, time.Date(2023, 12, 18, 0, 0, 0, 0, time.UTC))
		req, err := http.NewRequest(http.MethodGet, "/week", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		srv.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNotFound, rr.Code)
	})
}
