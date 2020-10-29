package data

import (
	"github.com/google/uuid"
	"net/url"
	"time"
)

type Feeds interface {
	GetFeedList() ([]Feed, error)
	GetFeedById(id uuid.UUID) (*Feed, error)
	AddFeed(name string, addr url.URL) (*Feed, error)
	DeleteFeed(id uuid.UUID) error
	UpdateFeed(feed Feed) error
}

type UnreadItems interface {
	AddItem(feed uuid.UUID, title string, link string, pub time.Time, desc string, content string) (*FeedItem, error)
	ListItems() ([]FeedItem, error)
	DeleteAllItems() error
}
