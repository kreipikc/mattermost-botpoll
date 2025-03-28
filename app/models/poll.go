package models

import "time"

type PollBody struct {
	Id          int
	Title       string
	Description string
	Variants    map[string]int
	DateEnd     time.Time
}
