package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func SerializeJSON(t testing.TB, v any) []byte {
	b, err := json.Marshal(v)
	require.NoError(t, err)

	return b
}

func SerializeJSONReader(t testing.TB, v any) io.Reader {
	b := SerializeJSON(t, v)

	return bytes.NewReader(b)
}
