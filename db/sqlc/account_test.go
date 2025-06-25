package db

import (
	"context"
	"testing"

	"github.com/akkahshh24/go-simple-bank-app/util"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	ctx := context.Background()
	arg := CreateAccountParams{
		Owner:    util.RandomName(),
		Balance:  util.RandomAmount(10000, 100000),
		Currency: util.RandomCurrency(),
	}

	account, err := testStore.CreateAccount(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, account.UpdatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	ctx := context.Background()
	createdAccount := createRandomAccount(t)
	gotAccount, err := testStore.GetAccount(ctx, createdAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gotAccount)
	require.Equal(t, createdAccount, gotAccount)
}

func TestUpdateAccount(t *testing.T) {
	ctx := context.Background()
	createdAccount := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      createdAccount.ID,
		Balance: util.RandomAmount(10000, 100000),
	}
	updatedAccount, err := testStore.UpdateAccount(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)
	require.Equal(t, createdAccount.ID, updatedAccount.ID)
	require.Equal(t, arg.Balance, updatedAccount.Balance)
}

func TestDeleteAccount(t *testing.T) {
	ctx := context.Background()
	createdAccount := createRandomAccount(t)
	err := testStore.DeleteAccount(ctx, createdAccount.ID)
	require.NoError(t, err)

	deletedAccount, err := testStore.GetAccount(ctx, createdAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, deletedAccount)
}

func TestListAccounts(t *testing.T) {
	ctx := context.Background()

	// Create multiple random accounts
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testStore.ListAccounts(ctx, arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
