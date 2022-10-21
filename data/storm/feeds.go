package storm

import (
	"github.com/asdine/storm/v3"
	"github.com/google/uuid"
	"net/url"
	"reader/data"
	"time"
)

type FeedsDatabase struct {
	db *storm.DB
}

func NewStormFeedsDatabase(path string) (*FeedsDatabase, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, err
	}

	feeds := &FeedsDatabase{
		db: db,
	}

	return feeds, nil
}

func (f *FeedsDatabase) GetFeedList() ([]data.Feed, error) {
	var feeds []data.Feed
	err := f.db.All(&feeds)
	if err != nil {
		return nil, err
	}

	return feeds, nil
}

func (f *FeedsDatabase) GetFeedById(id uuid.UUID) (*data.Feed, error) {
	var feed data.Feed
	err := f.db.One("ID", id, &feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}

func (f *FeedsDatabase) AddFeed(name string, addr url.URL) (*data.Feed, error) {
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

func (f *FeedsDatabase) DeleteFeed(id uuid.UUID) error {
	return f.db.DeleteStruct(&data.Feed{ID: id})
}

func (f *FeedsDatabase) UpdateFeed(feed data.Feed) error {
	//feed.AddTime = time.Time{} // can't remember why this is here

	err := f.db.Update(&feed)

	return err
}
