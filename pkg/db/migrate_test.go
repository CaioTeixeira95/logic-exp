package db

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"testing/fstest"

	"github.com/CaioTeixeira95/logic-exp/pkg/utils/testutils"
	"github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrate(t *testing.T) {
	fyleSystem := fstest.MapFS{
		"2023-04-16.0.create-foo-table.sql": {
			Data: []byte(`-- +migrate Up
CREATE TABLE public.foo (id SERIAL);
-- +migrate Down
DROP TABLE public.foo;
`),
		},
		"2023-04-16.1.create-bar-table.sql": {
			Data: []byte(`-- +migrate Up
CREATE TABLE public.bar (id SERIAL);
-- +migrate Down
DROP TABLE public.bar;
`),
		},
	}

	user := os.Getenv("PGUSER")
	if user == "" {
		user = "postgres:postgres"
	}

	dbURL := fmt.Sprintf("postgres://%s@localhost/?sslmode=disable", user)

	conn, err := Open(dbURL)
	require.NoError(t, err)

	testDBName := testutils.RandomTestDatabaseName()
	_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE %s", pq.QuoteIdentifier(testDBName)))
	require.NoError(t, err)

	dbTestURL := fmt.Sprintf("postgres://%s@localhost/%s?sslmode=disable", user, testDBName)

	n, err := Migrate(dbTestURL, migrate.Up, 2, http.FS(fyleSystem))
	require.NoError(t, err)

	assert.Equal(t, 2, n)

	_, err = conn.Exec(fmt.Sprintf("DROP DATABASE %s", pq.QuoteIdentifier(testDBName)))
	require.NoError(t, err)
}
