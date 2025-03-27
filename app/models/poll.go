package models

import "time"

type PollBody struct {
	Title       string
	Description string
	Variants    map[string]int
	DateEnd     time.Time
}
