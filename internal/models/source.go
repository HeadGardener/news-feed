package models

import (
	"errors"
	"time"
)

type Source struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	FeedURL   string    `db:"feed_url"`
	CreatedAt time.Time `db:"created_at"`
}

type SourceInput struct {
	Name      string `json:"name"`
	FeedURL   string `json:"feed_url"`
	CreatedAt time.Time
}

func (s *SourceInput) Validate() error {
	if s.FeedURL[:7] != "http://" && s.FeedURL[:8] != "https://" {
		return errors.New("invalid link")
	}

	s.CreatedAt = time.Now()
	return nil
}
