package repository

import (
	"fmt"
	"time"
)

type News struct {
	ID          int
	PubDate     time.Time
	Title       string
	Description string
	Link        string
}

func FormatSingleNews(n News) string {
	return fmt.Sprintf("Title: %s\nPublished Date: %s\nDescription: %s\nLink: %s\n",
		n.Title, n.PubDate.Format("2006-01-02 15:04:05"), n.Description, n.Link)
}
