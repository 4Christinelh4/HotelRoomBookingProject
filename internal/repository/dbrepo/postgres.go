package dbrepo

import (
	"context"
	"my/gomodule/internal/models"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `insert into reservations (
             email, phone, first_name, last_name, start_date, end_date, room_id
	) 
	values ($1, $2, $3, $4, $5, $6, $7) returning id`

	m.DB.QueryRowContext(ctx, stmt, res.Email, res.Phone, res.FirstName, res.LastName,
		time.Now(), time.Now(), res.RoomID).Scan(&newID)

	return newID, nil
}

func (m *postgresDBRepo) InsertRoomRestriction(res models.RoomRestrictions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (
             restriction_id, room_id, reservation_id, start_date, end_date, created_at, updated_at) 
				values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt, res.RestrictionID, res.RoomID, res.ReservationID,
		time.Now(), time.Now(), time.Now(), time.Now())

	if err != nil {
		return err
	}

	return nil
}
