package storm

import (
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"reader/data"
	"time"
)

type itemsDatabase struct {
	db *storm.DB
}

func NewStormItemsDatabase(path string) (*itemsDatabase, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, err
	}

	items := &itemsDatabase{
		db: db,
	}

	return items, nil
}

func (i *itemsDatabase) AddItem(feed uuid.UUID, title string, link string, pub time.Time, desc string, content string) (*data.FeedItem, error) {
	item := &data.FeedItem{
		ID:          uuid.New(),
		Feed:        feed,
		Title:       title,
		Link:        link,
		PubDate:     pub,
		Description: desc,
		Content:     content,
		Fetched:     time.Now(),
	}

	err := i.db.Save(item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (i *itemsDatabase) ListItems() ([]data.FeedItem, error) {
	var items []data.FeedItem
	err := i.db.AllByIndex("PubDate", &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (i *itemsDatabase) DeleteAllItems() error {
	return i.db.Select(q.True()).Delete(&data.FeedItem{})
}
