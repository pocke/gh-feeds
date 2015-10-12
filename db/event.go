//go:generate argen
package db

import (
	"time"

	"github.com/monochromegane/argen"
)

//+AR
type Event struct {
	ID          int `db:"pk"` // auto increment
	PublishedAt time.Time
	Type        string // WatchEvent, ...
	HTML        string
	UserId      int
	AuthorName  string
	URL         string
	ImageURL    string
}

func (m Event) belongsToUser() *ar.Association { return nil }
