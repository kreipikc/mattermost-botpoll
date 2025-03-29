package models

import "time"

type PollBody struct {
	Id          int
	AuthorID    string
	Title       string
	Description string
	Variants    map[string]int
	DateEnd     time.Time
}
