package dbrepo

import (
	"context"
	"my/gomodule/internal/models"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) InsertReservation(res models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into reservations (
             email, phone, first_name, last_name, start_date, end_date, room_id
	) 
	values ($1, $2, $3, $4, $5, $6, $7)`

	m.DB.ExecContext(ctx, stmt, res.Email, res.Phone, res.FirstName, res.LastName,
		time.Now(), time.Now(), res.RoomID)
	return nil
}
