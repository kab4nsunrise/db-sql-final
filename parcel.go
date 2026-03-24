package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)


type ParcelStore struct {
	db *sql.DB
}


func NewParcelStore(db *sql.DB) *ParcelStore {
	return &ParcelStore{db: db}
}


func (s *ParcelStore) Add(ctx context.Context, client int, address string) (int, error) {
	
	query := `INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)`
	createdAt := time.Now().UTC().Format(time.RFC3339) 
	status := ParcelStatusRegistered                   

	res, err := s.db.ExecContext(ctx, query, client, status, address, createdAt)
	if err != nil {
		return 0, fmt.Errorf("не удалось добавить посылку: %w", err)
	}

	
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("не удалось получить ID новой посылки: %w", err)
	}
	return int(id), nil
}


func (s *ParcelStore) GetByClient(ctx context.Context, client int) ([]Parcel, error) {
	query := `SELECT number, client, status, address, created_at FROM parcel WHERE client = ?`
	rows, err := s.db.QueryContext(ctx, query, client)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		parcels = append(parcels, p)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по строкам: %w", err)
	}
	return parcels, nil
}


func (s *ParcelStore) SetStatus(ctx context.Context, number int, status string) error {
	query := `UPDATE parcel SET status = ? WHERE number = ?`
	res, err := s.db.ExecContext(ctx, query, status, number)
	if err != nil {
		return fmt.Errorf("не удалось обновить статус: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество затронутых строк: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}


func (s *ParcelStore) SetAddress(ctx context.Context, number int, address string) error {
	query := `UPDATE parcel SET address = ? WHERE number = ?`
	res, err := s.db.ExecContext(ctx, query, address, number)
	if err != nil {
		return fmt.Errorf("не удалось обновить адрес: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество затронутых строк: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows 
	}
	return nil
}


func (s *ParcelStore) Delete(ctx context.Context, number int) error {
	query := `DELETE FROM parcel WHERE number = ?`
	res, err := s.db.ExecContext(ctx, query, number)
	if err != nil {
		return fmt.Errorf("не удалось удалить посылку: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество затронутых строк: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows 
	}
	return nil
}
