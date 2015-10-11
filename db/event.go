//go:generate argen
package db

import "time"

//+AR
type Event struct {
	ID          int `db:"pk"`
	PublishedAt time.Time
	Type        string // WatchEvent, ...
	HTML        string
	UserId      int
	AuthorName  string
	URL         string
	ImageURL    string
}
