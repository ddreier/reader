package data

import "net/url"

type Feeds interface {
	GetFeedList() ([]Feed, error)
	AddFeed(name string, addr url.URL) (*Feed, error)
}
