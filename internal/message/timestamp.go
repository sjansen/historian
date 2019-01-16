package message

import (
	"strconv"
	"time"
)

type Timestamp time.Time

func (t *Timestamp) Unix() int64 {
	return time.Time(*t).Unix()
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	x, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	*t = Timestamp(x)
	return nil
}
