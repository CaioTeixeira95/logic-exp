package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func Open(dataSourceName string) (*sql.DB, error) {
	return sql.Open("postgres", dataSourceName)
}
