package repository

import (
	"my/gomodule/internal/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestriction(res models.RoomRestrictions) error

	GetUserByID(id int) (models.User, error)

	UpdateUser(u models.User) error

	Authenticate(email, password string) (int, string, error)

	AllReservations() ([]models.Reservation, error)

	AllNewReservations() ([]models.Reservation, error)

	SearchAvailabilityByDates(roomID int, start, end time.Time) (bool, error)

	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)

	GetReservationByID(id int) (models.Reservation, error)

	UpdateReservation(u models.Reservation) error

	DeleteReservation(id int) error

	UpdatedProcessedForReservation(id, processed int) error
}
