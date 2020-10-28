package storm

import (
	"github.com/asdine/storm/v3"
	"github.com/google/uuid"
	"net/url"
	"reader/data"
	"time"
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

func (f *feedsDatabase) GetFeedById(id uuid.UUID) (*data.Feed, error) {
	var feed data.Feed
	err := f.db.One("ID", id, &feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}

func (f *feedsDatabase) AddFeed(name string, addr url.URL) (*data.Feed, error) {
	feed := data.Feed{
		ID:                uuid.New(),
		Name:              name,
		Addr:              addr,
		AddTime:           time.Now().UTC(),
		CheckTime:         time.Time{},
		MostRecentPubDate: time.Time{},
	}

	err := f.db.Save(&feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}

func (f *feedsDatabase) DeleteFeed(id uuid.UUID) error {
	return f.db.DeleteStruct(&data.Feed{ID: id})
}

func (f *feedsDatabase) UpdateFeed(feed data.Feed) error {
	feed.AddTime = time.Time{}

	err := f.db.Update(&feed)

	return err
}
