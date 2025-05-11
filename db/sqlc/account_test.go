package db

import (
	"context"
	"testing"

	"github.com/akkahshh24/go-simple-bank-app/util"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	ctx := context.Background()
	createdAccount := createRandomAccount(t)
	gotAccount, err := testQueries.GetAccount(ctx, createdAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gotAccount)
	require.Equal(t, createdAccount, gotAccount)
}

func TestUpdateAccount(t *testing.T) {
	ctx := context.Background()
	createdAccount := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      createdAccount.ID,
		Balance: util.RandomBalance(10000, 100000),
	}
	updatedAccount, err := testQueries.UpdateAccount(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)
	require.Equal(t, createdAccount.ID, updatedAccount.ID)
	require.Equal(t, arg.Balance, updatedAccount.Balance)
}

func TestDeleteAccount(t *testing.T) {
	ctx := context.Background()
	createdAccount := createRandomAccount(t)
	err := testQueries.DeleteAccount(ctx, createdAccount.ID)
	require.NoError(t, err)

	deletedAccount, err := testQueries.GetAccount(ctx, createdAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, deletedAccount)
}

func TestListAccounts(t *testing.T) {
	// Create multiple random accounts
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	ctx := context.Background()
	accounts, err := testQueries.ListAccounts(ctx, arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func createRandomAccount(t *testing.T) Account {
	ctx := context.Background()
	arg := CreateAccountParams{
		HolderName: util.RandomHolderName(),
		Balance:    util.RandomBalance(10000, 100000),
		Currency:   util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.HolderName, account.HolderName)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, account.UpdatedAt)

	return account
}
