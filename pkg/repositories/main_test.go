package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/CaioTeixeira95/logic-exp/migrations"
	"github.com/CaioTeixeira95/logic-exp/pkg/db"
	"github.com/CaioTeixeira95/logic-exp/pkg/utils/testutils"
	"github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

var testConn *sql.DB

func TestMain(m *testing.M) {
	user := os.Getenv("PGUSER")
	if user == "" {
		user = "postgres:postgres"
	}

	dbURL := fmt.Sprintf("postgres://%s@db/?sslmode=disable", user)

	conn, err := db.Open(dbURL)
	if err != nil {
		log.Fatalf("error connecting to the database: %s", err.Error())
	}
	defer conn.Close()

	testDBName := testutils.RandomTestDatabaseName()
	_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE %s", pq.QuoteIdentifier(testDBName)))
	if err != nil {
		log.Fatalf("error creating test database: %s", err.Error())
	}

	dbTestURL := fmt.Sprintf("postgres://%s@db/%s?sslmode=disable", user, testDBName)

	_, err = db.Migrate(dbTestURL, migrate.Up, 0, http.FS(migrations.FS))
	if err != nil {
		log.Fatalf("error applying migrations on test database: %s", err.Error())
	}

	testConn, err = db.Open(dbTestURL)
	if err != nil {
		log.Fatalf("error connection to database %s: %s", dbTestURL, err.Error())
	}

	exitCode := m.Run()

	testConn.Close()

	_, err = conn.Exec(fmt.Sprintf("DROP DATABASE %s", pq.QuoteIdentifier(testDBName)))
	if err != nil {
		log.Fatalf("error dropping test database: %s", err.Error())
	}

	os.Exit(exitCode)
}
