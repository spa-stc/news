package passwordutil_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"stpaulacademy.tech/newsletter/util/passwordutil"
)

func TestPasswordHashing(t *testing.T) {
	t.Parallel()

	password := "Hello World!"

	hash, err := passwordutil.HashPassword(password)
	require.NoError(t, err)

	ok, err := passwordutil.VerifyPassword(hash, password)
	require.NoError(t, err)
	require.True(t, ok)
}
