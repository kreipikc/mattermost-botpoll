package models

import "time"

type PollBody struct {
	Title       string
	Description string
	Variants    []string
	DateEnd     time.Time
}
