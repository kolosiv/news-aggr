package repository

import (
	"time"
)

type News struct {
	ID          int
	PubDate     time.Time
	Title       string
	Description string
	Link        string
	SourceName  string
}
