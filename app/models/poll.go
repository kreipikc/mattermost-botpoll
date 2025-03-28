package models

import "time"

type PollBody struct {
	Id          uint32
	AuthorID    string
	Title       string
	Description string
	Variants    map[string]int
	DateEnd     time.Time
}
