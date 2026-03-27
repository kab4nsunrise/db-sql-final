package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "modernc.org/sqlite"
)

const (
    ParcelStatusRegistered = "registered"
    ParcelStatusSent       = "sent"
    ParcelStatusDelivered  = "delivered"
)

func main() {
    
    db, err := sql.Open("sqlite", "tracker.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS parcel (
            number INTEGER PRIMARY KEY AUTOINCREMENT,
            client INTEGER,
            status TEXT,
            address TEXT,
            created_at TEXT
        )
    `)
    if err != nil {
        log.Fatal(err)
    }

    store := NewParcelStore(db)

    
    clientID := 1
    address := "ул. Пушкина, д. 10"
    p := Parcel{
        Client:  clientID,
        Status:  ParcelStatusRegistered,
        Address: address,
    }
    id, err := store.Add(p)
    if err != nil {
        log.Fatal("Ошибка при регистрации посылки:", err)
    }
    fmt.Printf("Посылка зарегистрирована, номер: %d\n", id)

    
    newAddress := "ул. Лермонтова, д. 5"
    err = store.SetAddress(id, newAddress)
    if err != nil {
        log.Fatal("Ошибка при изменении адреса:", err)
    }
    fmt.Printf("Адрес посылки %d изменён на %s\n", id, newAddress)

    
    err = store.SetStatus(id, ParcelStatusSent)
    if err != nil {
        log.Fatal("Ошибка при изменении статуса:", err)
    }
    fmt.Printf("Статус посылки %d изменён на %s\n", id, ParcelStatusSent)

    
    parcels, err := store.GetByClient(clientID)
    if err != nil {
        log.Fatal("Ошибка при получении списка посылок:", err)
    }
    fmt.Printf("Посылки клиента %d:\n", clientID)
    for _, p := range parcels {
        fmt.Printf("  Номер: %d, Статус: %s, Адрес: %s\n", p.Number, p.Status, p.Address)
    }

    
    err = store.Delete(id)
    if err != nil {
        log.Fatal("Ошибка при удалении посылки:", err)
    }
    fmt.Printf("Посылка %d удалена\n", id)
}
