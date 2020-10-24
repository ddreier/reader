package data

import (
	"github.com/google/uuid"
	"net/url"
)

type Feed struct {
	ID   uuid.UUID
	Name string
	Addr url.URL
}
