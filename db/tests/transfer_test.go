package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/BogoCvetkov/go_mastercalss/db"
	gen "github.com/BogoCvetkov/go_mastercalss/db/generated"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, account1, account2 gen.Account) gen.Transfer {
	arg := gen.CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        RandomMoney(),
	}

	transfer, err := testStore.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, account1, account2)

	transfer2, err := testStore.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, account1, account2)
		createRandomTransfer(t, account2, account1)
	}

	arg := gen.ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testStore.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}
}

func TestTransferTrxConcurrent(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	fmt.Println(">> before:", acc1.Balance, acc2.Balance)

	amount := int64(10)

	errchan := make(chan error)
	reschan := make(chan *db.TransferTxResult)

	// n concurrent transfer transactions
	n := 10

	for i := 0; i < n; i++ {
		data := gen.CreateTransferParams{
			FromAccountID: acc1.ID,
			ToAccountID:   acc2.ID,
			Amount:        amount,
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			res, err := testStore.TransferTrx(ctx, data)

			errchan <- err
			reschan <- res
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errchan
		require.NoError(t, err)

		result := <-reschan
		fmt.Printf("%d Transfer -->%v -- %v", i, result.FromAccount, result.ToAccount)
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, acc1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, acc2.ID, toAccount.ID)

		// check balances
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)

	}

	// check the final updated balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, acc1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, acc2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferTrxDeadlock(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	fmt.Println(">> before:", acc1.Balance, acc2.Balance)

	amount := int64(10)

	errchan := make(chan error)
	reschan := make(chan *db.TransferTxResult)

	// n concurrent transfer transactions
	n := 10

	for i := 0; i < n; i++ {
		data := gen.CreateTransferParams{
			Amount: amount,
		}

		//rotate accounts
		if i%2 == 0 {
			data.FromAccountID = acc1.ID
			data.ToAccountID = acc2.ID
		} else {
			data.FromAccountID = acc2.ID
			data.ToAccountID = acc1.ID
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			res, err := testStore.TransferTrx(ctx, data)

			errchan <- err
			reschan <- res
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errchan
		require.NoError(t, err)

	}

	// check the final updated balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, acc1.Balance, updatedAccount1.Balance)
	require.Equal(t, acc2.Balance, updatedAccount2.Balance)
}
