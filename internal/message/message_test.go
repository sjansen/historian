package message

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	require := require.New(t)

	r, err := os.Open("testdata/message.json")
	require.NoError(err)

	changed, err := time.Parse(time.RFC3339, "2017-06-27T16:20:44Z")
	require.NoError(err)

	expected := &Message{
		ID:        "595285dc-9c43-4b9c-a1e6-0cd9aff5b084",
		Version:   "v1",
		MemberID:  "56d8a839-1c52-437f-b981-c3a15a11d6d4",
		PrimaryID: float64(16927),
		ChangedAt: Timestamp(changed),
	}
	actual, err := Parse(r)
	require.NoError(err)
	require.Equal(expected, actual)
	require.Equal(changed.Unix(), actual.ChangedAt.Unix())
}
