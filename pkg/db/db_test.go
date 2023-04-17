package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	user := os.Getenv("PGUSER")
	if user == "" {
		user = "postgres:postgres"
	}

	conn, err := Open(fmt.Sprintf("postgres://%s@localhost/?sslmode=disable", user))
	require.NoError(t, err)

	var res int64
	err = conn.QueryRow("SELECT 1").Scan(&res)
	require.NoError(t, err)

	assert.Equal(t, int64(1), res)
}
