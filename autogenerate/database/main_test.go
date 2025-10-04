package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/WhiteMaks/go_academy_portal_service/util"

	_ "github.com/lib/pq"
)

var testDatabase *sql.DB
var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	testDatabase, err = sql.Open("postgres", config.Database.ConnectionString)
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}

	testQueries = New(testDatabase)

	os.Exit(m.Run())
}
