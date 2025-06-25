package db

import (
	"context"
	"testing"

	"github.com/akkahshh24/go-simple-bank-app/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, fromAccount Account, toAccount Account) Transfer {
	ctx := context.Background()
	arg := CreateTransferParams{
		FromAccountID: util.Int64ToPgTypeInt8(fromAccount.ID),
		ToAccountID:   util.Int64ToPgTypeInt8(toAccount.ID),
		Amount:        util.RandomAmount(100, 1000),
	}
	transfer, err := testStore.CreateTransfer(ctx, arg)
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
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	createRandomTransfer(t, fromAccount, toAccount)
}

func TestGetTransfer(t *testing.T) {
	ctx := context.Background()
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	createdTransfer := createRandomTransfer(t, fromAccount, toAccount)

	gotTransfer, err := testStore.GetTransfer(ctx, createdTransfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gotTransfer)

	require.Equal(t, createdTransfer.ID, gotTransfer.ID)
	require.Equal(t, createdTransfer.FromAccountID, gotTransfer.FromAccountID)
	require.Equal(t, createdTransfer.ToAccountID, gotTransfer.ToAccountID)
	require.Equal(t, createdTransfer.Amount, gotTransfer.Amount)
	require.Equal(t, createdTransfer.CreatedAt.Time, gotTransfer.CreatedAt.Time)
}

func TestListTransfers(t *testing.T) {
	ctx := context.Background()
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	// Create multiple random transfers
	for range 10 {
		createRandomTransfer(t, fromAccount, toAccount)
	}

	arg := ListTransfersParams{
		FromAccountID: util.Int64ToPgTypeInt8(fromAccount.ID),
		ToAccountID:   util.Int64ToPgTypeInt8(toAccount.ID),
		Limit:         5,
		Offset:        5,
	}
	transfers, err := testStore.ListTransfers(ctx, arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID.Int64 == fromAccount.ID || transfer.ToAccountID.Int64 == toAccount.ID)
	}
}
