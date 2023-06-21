package db

import (
	"database/sql"

	db "github.com/BogoCvetkov/go_mastercalss/db/generated"
)

type Store struct {
	*db.Queries
	conn *sql.DB
}

func NewStore(conn *sql.DB) *Store {

	return &Store{
		conn:    conn,
		Queries: db.New(conn),
	}
}
