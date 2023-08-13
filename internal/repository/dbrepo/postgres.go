package dbrepo

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
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

// GetUserByID searches the user by ID
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, email, first_name, last_name, password, access_level from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	var u models.User
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Password, &u.AccessLevel)

	if err != nil {
		return u, err
	}

	return u, nil
}

func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name=$1, last_name=$2, email=$3, access_level=$5`
	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.AccessLevel)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) Authenticate(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashPw string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email=$1", email)
	err := row.Scan(&id, &hashPw)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPw), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashPw, nil
}

func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation
	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, 
		       r.end_date, r.room_id, r.processed,
		        rm.id, rm.room_name
				from reservations r left join rooms rm on rm.id = r.room_id
		     	order by r.start_date asc
	`

	rows, err := m.DB.QueryContext(ctx, query)
	defer rows.Close()
	if err != nil {
		return reservations, err
	}

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(&i.ID, &i.FirstName, &i.LastName,
			&i.Email, &i.Phone, &i.StartDate, &i.EndDate, &i.RoomID, &i.Processed, &i.Room.ID, &i.Room.RoomName)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// SearchAvailabilityByDates: check if there is any overlaps
func (m *postgresDBRepo) SearchAvailabilityByDates(roomID int, start, end time.Time) (bool, error) {
	//TODO implement me
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select count(id) from room_restrictions 
                 where room_id = $1 and $2 < end_date  and $3 > start_date)`

	var numRows int
	searchResults := m.DB.QueryRowContext(ctx, query, roomID, start, end)

	err := searchResults.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `select r.id, r.room_name from rooms r
                 where r.id not in 
                       (select room_id from room_restrictions rr
                        where rr.start_date < $2 or rr.end_date > $1)`

	searchResults, err := m.DB.QueryContext(ctx, query, start, end)

	if err != nil {
		return rooms, err
	}

	for searchResults.Next() {
		var room models.Room
		searchResults.Scan(&room.ID, &room.RoomName)

		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation
	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, 
		       r.room_id, r.processed,
		        rm.id, rm.room_name
				from reservations r left join rooms rm on rm.id = r.room_id
				where r.processed= 0 
		     	order by r.start_date asc
	`

	rows, err := m.DB.QueryContext(ctx, query)
	defer rows.Close()
	if err != nil {
		return reservations, err
	}

	for rows.Next() {
		var i models.Reservation

		err := rows.Scan(&i.ID, &i.FirstName, &i.LastName,
			&i.Email, &i.Phone, &i.StartDate, &i.EndDate, &i.RoomID, &i.Processed,
			&i.Room.ID, &i.Room.RoomName)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservation
	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, 
				r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
				rm.id, rm.room_name from reservations r left join rooms rm on (r.room_id = rm.id)
				where r.id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(&res.ID, &res.FirstName, &res.LastName,
		&res.Email, &res.Phone, &res.StartDate, &res.EndDate, &res.RoomID,
		&res.CreatedAt, &res.UpdatedAt, &res.Processed,
		&res.Room.ID, &res.Room.RoomName)

	if err != nil {
		return res, err
	}
	return res, nil
}

func (m *postgresDBRepo) UpdateReservation(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set first_name=$1, last_name=$2, email=$3, 
                        phone=$4, updated_at=$5 where id=$6`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName, u.LastName, u.Email, u.Phone, u.UpdatedAt, u.ID)

	return err
}

func (m *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from reservations where id=$1`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m *postgresDBRepo) UpdatedProcessedForReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set processed=$1 where id=$2`

	_, err := m.DB.ExecContext(ctx, query, processed, id)
	return err
}
