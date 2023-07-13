package db

import (
	"context"
	"errors"
	"fmt"

	db "github.com/BogoCvetkov/go_mastercalss/db/generated"
)

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    db.Transfer `json:"transfer"`
	FromAccount db.Account  `json:"from_account"`
	ToAccount   db.Account  `json:"to_account"`
	FromEntry   db.Entry    `json:"from_entry"`
	ToEntry     db.Entry    `json:"to_entry"`
}

func (s *Store) TransferTrx(ctx context.Context, data db.CreateTransferParams) (*TransferTxResult, error) {

	var result TransferTxResult

	// Begin Trx --->
	tx, err := s.conn.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	qtx := s.WithTx(tx)

	// Check if sender has enough balance
	from, err := qtx.GetAccount(ctx, data.FromAccountID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if from.Balance < data.Amount {
		return nil, errors.New("not enough balance for transfer")
	}

	// Save transfer info
	result.Transfer, err = qtx.CreateTransfer(ctx, data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Create entry for the OUTGOING money
	e1 := db.CreateEntryParams{
		AccountID: data.FromAccountID,
		Amount:    -data.Amount,
	}

	result.FromEntry, err = qtx.CreateEntry(ctx, e1)
	if err != nil {
		return nil, err
	}

	// Create entry for the INCOMING money
	e2 := db.CreateEntryParams{
		AccountID: data.ToAccountID,
		Amount:    data.Amount,
	}
	result.ToEntry, err = qtx.CreateEntry(ctx, e2)
	if err != nil {
		return nil, err
	}

	// Update Balances
	if data.FromAccountID < data.ToAccountID {
		r1, r2, _ := moveMoney(ctx, data.FromAccountID, data.ToAccountID, -data.Amount, data.Amount, qtx, &result)
		result.FromAccount, result.ToAccount = *r1, *r2
	} else {
		r1, r2, _ := moveMoney(ctx, data.ToAccountID, data.FromAccountID, data.Amount, -data.Amount, qtx, &result)
		result.ToAccount, result.FromAccount = *r1, *r2
	}

	// Commit Trx <---
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func moveMoney(ctx context.Context, a1 int64, a2 int64, m1 int64, m2 int64, qtx *db.Queries, r *TransferTxResult) (*db.Account, *db.Account, error) {
	var err error
	var acc1 db.Account
	var acc2 db.Account

	// Acc --> 1
	b1 := db.AddAccountBalanceParams{
		// substract
		Amount: m1,
		ID:     a1,
	}
	acc1, err = qtx.AddAccountBalance(ctx, b1)
	if err != nil {
		return nil, nil, err
	}

	// Acc --> 2
	b2 := db.AddAccountBalanceParams{
		// increment
		Amount: m2,
		ID:     a2,
	}
	acc2, err = qtx.AddAccountBalance(ctx, b2)
	if err != nil {
		return nil, nil, err
	}

	return &acc1, &acc2, nil
}
