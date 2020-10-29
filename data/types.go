package data

import (
	"github.com/google/uuid"
	"net/url"
	"time"
)

type Feed struct {
	ID                uuid.UUID
	Name              string
	Addr              url.URL
	AddTime           time.Time
	CheckTime         time.Time
	MostRecentPubDate time.Time
}

type FeedItem struct {
	ID          uuid.UUID
	Feed        uuid.UUID
	Title       string
	Link        string
	PubDate     time.Time `storm:"index"`
	Description string
	Content     string
	Fetched     time.Time
}
