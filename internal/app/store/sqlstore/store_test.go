package sqlstore_test

import (
	"os"
	"testing"
)

var databaseURL string

// TestMain method for initializing db url for testing db from env
func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost dbname=auth_db_test user=user password=password sslmode=disable"
	}

	os.Exit(m.Run())
}
