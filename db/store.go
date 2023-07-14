package db

import (
	"context"
	"database/sql"

	db "github.com/BogoCvetkov/go_mastercalss/db/generated"
	"github.com/hibiken/asynq"
)

// Store defines all functions to execute db queries and transactions
type IStore interface {
	db.Querier
	TransferTrx(ctx context.Context, arg db.CreateTransferParams) (*TransferTxResult, error)
	CreateUserTrx(ctx context.Context, arg db.CreateUserParams, asq *asynq.Client) (*db.User, error)
}

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
