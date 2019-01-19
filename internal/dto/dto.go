package dto

import "time"

type Message struct {
	Timestamp time.Time
	RawData   string
}
