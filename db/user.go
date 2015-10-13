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

func CreateUser(p *UserParams) (*User, error) {
	ins := ar.NewInsert(db, logger).Table("users").Params(map[string]interface{}{
		"id":   p.ID,
		"name": p.Name,
		"auth": p.Auth,
	})
	if _, err := ins.Exec(); err != nil {
		errs := &ar.Errors{}
		errs.AddError("base", err)
		return nil, errs
	}
	return &User{
		ID:   p.ID,
		Name: p.Name,
		Auth: p.Auth,
	}, nil
}
