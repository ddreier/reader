package storm

import (
	"github.com/asdine/storm/v3"
	"github.com/google/uuid"
	"net/url"
	"reader/data"
)

type feedsDatabase struct {
	db *storm.DB
}

func NewStormFeedsDatabase(path string) (*feedsDatabase, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, err
	}

	feeds := &feedsDatabase{
		db: db,
	}

	return feeds, nil
}

func (f *feedsDatabase) GetFeedList() ([]data.Feed, error) {
	var feeds []data.Feed
	err := f.db.All(&feeds)
	if err != nil {
		return nil, err
	}

	return feeds, nil
}

func (f *feedsDatabase) AddFeed(name string, addr url.URL) (*data.Feed, error) {
	feed := data.Feed{
		ID:   uuid.New(),
		Name: name,
		Addr: addr,
	}

	err := f.db.Save(&feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}
