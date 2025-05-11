package db

import (
	"context"
	"testing"

	"github.com/akkahshh24/go-simple-bank-app/util"
	"github.com/stretchr/testify/require"
)

func TestCreateTransactions(t *testing.T) {
	// Create two random accounts for testing transfer
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransaction(t, account1, account2)
}

func TestGetTransaction(t *testing.T) {
	ctx := context.Background()
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createdTransaction := createRandomTransaction(t, account1, account2)

	gotTransaction, err := testQueries.GetTransaction(ctx, createdTransaction.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gotTransaction)
	require.Equal(t, createdTransaction.ID, gotTransaction.ID)
	require.Equal(t, createdTransaction.FromAccountID, gotTransaction.FromAccountID)
	require.Equal(t, createdTransaction.ToAccountID, gotTransaction.ToAccountID)
	require.Equal(t, createdTransaction.Type, gotTransaction.Type)
	require.Equal(t, createdTransaction.Amount, gotTransaction.Amount)
	require.Equal(t, createdTransaction.Description, gotTransaction.Description)
	require.Equal(t, createdTransaction.CreatedAt, gotTransaction.CreatedAt)
}

func TestListTransactions(t *testing.T) {
	ctx := context.Background()
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Create multiple random transactions
	for i := 0; i < 5; i++ {
		createRandomTransaction(t, account1, account2)
		createRandomTransaction(t, account2, account1)
	}

	arg := ListTransactionsParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Limit:         5,
		Offset:        5,
	}

	transactions, err := testQueries.ListTransactions(ctx, arg)
	require.NoError(t, err)
	require.Len(t, transactions, 5)

	for _, transaction := range transactions {
		require.NotEmpty(t, transaction)
	}
}

func createRandomTransaction(t *testing.T, account1, account2 Account) Transaction {
	ctx := context.Background()
	arg := CreateTransactionParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Type:          util.RandomType(),
		Amount:        util.RandomAmount(10000, 1000000),
		Description:   util.RandomDescription(3),
	}

	transaction, err := testQueries.CreateTransaction(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	require.Equal(t, arg.FromAccountID, transaction.FromAccountID)
	require.Equal(t, arg.ToAccountID, transaction.ToAccountID)
	require.Equal(t, arg.Type, transaction.Type)
	require.Equal(t, arg.Amount, transaction.Amount)
	require.Equal(t, arg.Description, transaction.Description)

	require.NotZero(t, transaction.ID)
	require.NotZero(t, transaction.CreatedAt)
	return transaction
}
