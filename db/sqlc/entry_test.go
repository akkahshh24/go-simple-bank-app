package db

import (
	"context"
	"testing"

	"github.com/akkahshh24/go-simple-bank-app/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	ctx := context.Background()
	arg := CreateEntryParams{
		AccountID: util.Int64ToPgTypeInt8(account.ID),
		Amount:    util.RandomAmount(100, 1000),
	}

	entry, err := testStore.CreateEntry(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	ctx := context.Background()
	account := createRandomAccount(t)
	createdEntry := createRandomEntry(t, account)
	gotEntry, err := testStore.GetEntry(ctx, createdEntry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gotEntry)
	require.Equal(t, createdEntry.ID, gotEntry.ID)
	require.Equal(t, createdEntry.AccountID, gotEntry.AccountID)
	require.Equal(t, createdEntry.Amount, gotEntry.Amount)
	require.Equal(t, createdEntry.CreatedAt.Time, gotEntry.CreatedAt.Time)
}

func TestListEntries(t *testing.T) {
	ctx := context.Background()
	account := createRandomAccount(t)

	// Create multiple random entries
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	arg := ListEntriesParams{
		AccountID: util.Int64ToPgTypeInt8(account.ID),
		Limit:     5,
		Offset:    5,
	}
	entries, err := testStore.ListEntries(ctx, arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, account.ID, entry.AccountID.Int64)
	}
}
