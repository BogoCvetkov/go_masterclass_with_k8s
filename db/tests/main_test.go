package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	db "github.com/BogoCvetkov/go_mastercalss/db"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbUrl    = "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable"
)

var testStore *db.Store

func TestMain(m *testing.M) {

	conn, err := sql.Open(dbDriver, dbUrl)

	if err != nil {
		log.Fatal("Failed connecting to DB", err)
	}

	testStore = db.NewStore(conn)

	os.Exit(m.Run())
}
