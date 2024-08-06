package sliceutil_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"stpaulacademy.tech/newsletter/util/sliceutil"
)

func TestMap(t *testing.T) {
	t.Parallel()

	result1 := sliceutil.Map([]int{1, 2, 3, 4}, func(_ int) string {
		return "Hello"
	})

	result2 := sliceutil.Map([]int64{1, 2, 3, 4}, func(x int64) string {
		return strconv.FormatInt(x, 10)
	})

	require.Len(t, result1, 4)
	require.Len(t, result2, 4)
	require.Equal(t, []string{"Hello", "Hello", "Hello", "Hello"}, result1)
	require.Equal(t, []string{"1", "2", "3", "4"}, result2)
}
