package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func UseProd() error {
	d, err := sql.Open("mysql", "root:@/ghfeeds")
	if err != nil {
		return err
	}
	Use(d)
	return nil
}
