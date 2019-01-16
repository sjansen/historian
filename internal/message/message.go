package message

import (
	"encoding/json"
	"io"
)

type Message struct {
	ID        string      `json:"id"`
	Version   string      `json:"version"`
	MemberID  string      `json:"member_id"`
	PrimaryID interface{} `json:"primary_id"`
	ChangedAt Timestamp   `json:"changed_at"`
}

func Parse(r io.Reader) (*Message, error) {
	dec := json.NewDecoder(r)

	m := &Message{}
	if err := dec.Decode(m); err != nil {
		return nil, err
	}

	return m, nil
}
