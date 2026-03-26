package main

import (
	"context"
	"database/sql"
	"math/rand"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS parcel (
			number INTEGER PRIMARY KEY AUTOINCREMENT,
			client INTEGER,
			status TEXT,
			address TEXT,
			created_at TEXT
		)
	`)
	require.NoError(t, err)

	return db
}

func TestAddGetDelete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()
	ctx := context.Background()

	id, err := store.Add(ctx, parcel.Client, parcel.Address)
	require.NoError(t, err)
	require.NotZero(t, id)

	parcels, err := store.GetByClient(ctx, parcel.Client)
	require.NoError(t, err)
	require.Len(t, parcels, 1)

	stored := parcels[0]
	require.Equal(t, id, stored.Number)
	require.Equal(t, parcel.Client, stored.Client)
	require.Equal(t, ParcelStatusRegistered, stored.Status)
	require.Equal(t, parcel.Address, stored.Address)
	require.NotEmpty(t, stored.CreatedAt)

	err = store.Delete(ctx, id)
	require.NoError(t, err)

	parcels, err = store.GetByClient(ctx, parcel.Client)
	require.NoError(t, err)
	require.Len(t, parcels, 0)
}

func TestSetAddress(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()
	ctx := context.Background()

	id, err := store.Add(ctx, parcel.Client, parcel.Address)
	require.NoError(t, err)

	newAddress := "new test address"
	err = store.SetAddress(ctx, id, newAddress)
	require.NoError(t, err)

	parcels, err := store.GetByClient(ctx, parcel.Client)
	require.NoError(t, err)
	require.Len(t, parcels, 1)
	require.Equal(t, newAddress, parcels[0].Address)
}

func TestSetStatus(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()
	ctx := context.Background()

	id, err := store.Add(ctx, parcel.Client, parcel.Address)
	require.NoError(t, err)

	err = store.SetStatus(ctx, id, ParcelStatusSent)
	require.NoError(t, err)

	parcels, err := store.GetByClient(ctx, parcel.Client)
	require.NoError(t, err)
	require.Len(t, parcels, 1)
	require.Equal(t, ParcelStatusSent, parcels[0].Status)
}

func TestGetByClient(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewParcelStore(db)
	ctx := context.Background()

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := make(map[int]Parcel)

	client := randRange.Intn(10_000_000)
	for i := range parcels {
		parcels[i].Client = client
	}

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(ctx, parcels[i].Client, parcels[i].Address)
		require.NoError(t, err)
		require.NotZero(t, id)

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(ctx, client)
	require.NoError(t, err)
	require.Len(t, storedParcels, len(parcels))

	for _, stored := range storedParcels {
		original, ok := parcelMap[stored.Number]
		require.True(t, ok, "посылка с номером %d не найдена", stored.Number)
		require.Equal(t, original.Client, stored.Client)
		require.Equal(t, original.Address, stored.Address)
		require.Equal(t, ParcelStatusRegistered, stored.Status)
		require.NotEmpty(t, stored.CreatedAt)
	}
}
