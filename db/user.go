//go:generate argen
package db

import "github.com/monochromegane/argen"

//+AR
type User struct {
	ID   int `db:"pk"`
	Name string
	Auth string // ?
}

func (m User) hasManyEvents() *ar.Association { return nil }
