package main

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestAddGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE parcel (
			number INTEGER PRIMARY KEY AUTOINCREMENT,
			client INTEGER,
			status TEXT,
			address TEXT,
			created_at TEXT
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	store := NewParcelStore(db)

	p1 := Parcel{Client: 1, Status: "registered", Address: "addr1", CreatedAt: "now"}
	id1, err := store.Add(p1)
	if err != nil || id1 == 0 {
		t.Fatalf("Add failed: %v, id=%d", err, id1)
	}

	p2 := Parcel{Client: 1, Status: "sent", Address: "addr2", CreatedAt: "now"}
	id2, err := store.Add(p2)
	if err != nil || id2 == 0 {
		t.Fatalf("Add failed: %v, id=%d", err, id2)
	}

	p3 := Parcel{Client: 2, Status: "registered", Address: "addr3", CreatedAt: "now"}
	id3, err := store.Add(p3)
	if err != nil || id3 == 0 {
		t.Fatalf("Add failed: %v, id=%d", err, id3)
	}

	parcels, err := store.GetByClient(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(parcels) != 2 {
		t.Errorf("Expected 2 parcels, got %d", len(parcels))
	}

	found1, found2 := false, false
	for _, p := range parcels {
		if p.Number == id1 {
			if p.Address != "addr1" || p.Status != "registered" {
				t.Errorf("Parcel %d mismatch: %+v", id1, p)
			}
			found1 = true
		}
		if p.Number == id2 {
			if p.Address != "addr2" || p.Status != "sent" {
				t.Errorf("Parcel %d mismatch: %+v", id2, p)
			}
			found2 = true
		}
	}
	if !found1 || !found2 {
		t.Error("Not all parcels found for client 1")
	}

	for _, p := range parcels {
		if p.Number == id3 {
			t.Error("Client 2 parcel returned for client 1")
		}
	}
}

func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE parcel (
			number INTEGER PRIMARY KEY AUTOINCREMENT,
			client INTEGER,
			status TEXT,
			address TEXT,
			created_at TEXT
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	store := NewParcelStore(db)

	p := Parcel{Client: 1, Status: "registered", Address: "addr", CreatedAt: "now"}
	id, err := store.Add(p)
	if err != nil {
		t.Fatal(err)
	}

	err = store.SetStatus(id, "sent")
	if err != nil {
		t.Fatal(err)
	}

	parcels, err := store.GetByClient(1)
	if err != nil || len(parcels) != 1 {
		t.Fatalf("GetByClient failed: %v, len=%d", err, len(parcels))
	}
	if parcels[0].Status != "sent" {
		t.Errorf("Expected status 'sent', got '%s'", parcels[0].Status)
	}
}

func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE parcel (
			number INTEGER PRIMARY KEY AUTOINCREMENT,
			client INTEGER,
			status TEXT,
			address TEXT,
			created_at TEXT
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	store := NewParcelStore(db)

	p := Parcel{Client: 1, Status: "registered", Address: "old", CreatedAt: "now"}
	id, err := store.Add(p)
	if err != nil {
		t.Fatal(err)
	}

	newAddr := "new"
	err = store.SetAddress(id, newAddr)
	if err != nil {
		t.Fatal(err)
	}

	parcels, err := store.GetByClient(1)
	if err != nil || len(parcels) != 1 {
		t.Fatalf("GetByClient failed: %v, len=%d", err, len(parcels))
	}
	if parcels[0].Address != newAddr {
		t.Errorf("Expected address '%s', got '%s'", newAddr, parcels[0].Address)
	}
}

func TestDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE parcel (
			number INTEGER PRIMARY KEY AUTOINCREMENT,
			client INTEGER,
			status TEXT,
			address TEXT,
			created_at TEXT
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	store := NewParcelStore(db)

	p := Parcel{Client: 1, Status: "registered", Address: "addr", CreatedAt: "now"}
	id, err := store.Add(p)
	if err != nil {
		t.Fatal(err)
	}

	err = store.Delete(id)
	if err != nil {
		t.Fatal(err)
	}

	parcels, err := store.GetByClient(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(parcels) != 0 {
		t.Errorf("Expected 0 parcels after delete, got %d", len(parcels))
	}
}
