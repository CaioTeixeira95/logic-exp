package db

import (
	"fmt"
	"net/http"

	migrate "github.com/rubenv/sql-migrate"
)

func Migrate(dataSourceName string, dir migrate.MigrationDirection, count int, fileSystem http.FileSystem) (int, error) {
	conn, err := Open(dataSourceName)
	if err != nil {
		return 0, fmt.Errorf("database URL '%s': %w", dataSourceName, err)
	}
	defer conn.Close()

	ms := migrate.MigrationSet{}

	m := migrate.HttpFileSystemMigrationSource{FileSystem: fileSystem}
	return ms.ExecMax(conn, "postgres", m, dir, count)
}
