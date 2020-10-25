package data

import (
	"github.com/google/uuid"
	"net/url"
)

type Feeds interface {
	GetFeedList() ([]Feed, error)
	GetFeedById(id uuid.UUID) (*Feed, error)
	AddFeed(name string, addr url.URL) (*Feed, error)
	DeleteFeed(id uuid.UUID) error
	UpdateFeed(feed Feed) error
}

type FeedItems interface {
}
