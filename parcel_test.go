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

    p1 := Parcel{Client: 1, Status: "registered", Address: "ул. Пушкина, д. 1"}
    id1, err := store.Add(p1)
    if err != nil {
        t.Fatalf("Add failed: %v", err)
    }
    if id1 == 0 {
        t.Error("Add returned zero id")
    }

    p2 := Parcel{Client: 1, Status: "sent", Address: "ул. Лермонтова, д. 2"}
    id2, err := store.Add(p2)
    if err != nil {
        t.Fatalf("Add failed: %v", err)
    }
    if id2 == 0 {
        t.Error("Add returned zero id")
    }

    p3 := Parcel{Client: 2, Status: "registered", Address: "пр. Мира, д. 3"}
    id3, err := store.Add(p3)
    if err != nil {
        t.Fatalf("Add failed: %v", err)
    }
    if id3 == 0 {
        t.Error("Add returned zero id")
    }

    parcels, err := store.GetByClient(1)
    if err != nil {
        t.Fatalf("GetByClient failed: %v", err)
    }
    if len(parcels) != 2 {
        t.Errorf("Expected 2 parcels, got %d", len(parcels))
    }

    found1, found2 := false, false
    for _, p := range parcels {
        if p.Number == id1 {
            if p.Status != "registered" || p.Address != "ул. Пушкина, д. 1" || p.Client != 1 {
                t.Errorf("Parcel %d mismatch: %+v", id1, p)
            }
            found1 = true
        }
        if p.Number == id2 {
            if p.Status != "sent" || p.Address != "ул. Лермонтова, д. 2" || p.Client != 1 {
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

    p := Parcel{Client: 1, Status: "registered", Address: "ул. Пушкина, д. 1"}
    id, err := store.Add(p)
    if err != nil {
        t.Fatalf("Add failed: %v", err)
    }

    err = store.SetStatus(id, "sent")
    if err != nil {
        t.Fatalf("SetStatus failed: %v", err)
    }

    parcels, err := store.GetByClient(1)
    if err != nil {
        t.Fatalf("GetByClient failed: %v", err)
    }
    if len(parcels) != 1 {
        t.Fatalf("Expected 1 parcel, got %d", len(parcels))
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

    p := Parcel{Client: 1, Status: "registered", Address: "ул. Пушкина, д. 1"}
    id, err := store.Add(p)
    if err != nil {
        t.Fatalf("Add failed: %v", err)
    }

    newAddress := "ул. Лермонтова, д. 2"
    err = store.SetAddress(id, newAddress)
    if err != nil {
        t.Fatalf("SetAddress failed: %v", err)
    }

    parcels, err := store.GetByClient(1)
    if err != nil {
        t.Fatalf("GetByClient failed: %v", err)
    }
    if len(parcels) != 1 {
        t.Fatalf("Expected 1 parcel, got %d", len(parcels))
    }
    if parcels[0].Address != newAddress {
        t.Errorf("Expected address '%s', got '%s'", newAddress, parcels[0].Address)
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

    p := Parcel{Client: 1, Status: "registered", Address: "ул. Пушкина, д. 1"}
    id, err := store.Add(p)
    if err != nil {
        t.Fatalf("Add failed: %v", err)
    }

    err = store.Delete(id)
    if err != nil {
        t.Fatalf("Delete failed: %v", err)
    }

    parcels, err := store.GetByClient(1)
    if err != nil {
        t.Fatalf("GetByClient failed: %v", err)
    }
    if len(parcels) != 0 {
        t.Errorf("Expected 0 parcels after delete, got %d", len(parcels))
    }
}
