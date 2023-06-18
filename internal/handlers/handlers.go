package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"my/gomodule/internal/config"
	"my/gomodule/internal/driver"
	"my/gomodule/internal/forms"
	"my/gomodule/internal/helpers"
	"my/gomodule/internal/models"
	"my/gomodule/internal/renders"
	"my/gomodule/internal/repository"
	"my/gomodule/internal/repository/dbrepo"
	"net/http"
	"strconv"
	"time"
)

// Repo the repository used by the handlers
var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// Home: home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

type jsonResp struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {
	//start := r.Form.Get("started")
	//end := r.Form.Get("end")
	//w.Write([]byte(fmt.Sprintf("Posted method, start = %s, end = %s", start, end)))
	resp := jsonResp{
		OK:      true,
		Message: "available",
	}

	out, err := json.MarshalIndent(resp, "", "	")
	if err != nil {
		log.Print(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get error from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})

	data["reservation"] = reservation
	renders.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// Reservation handles reservation form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	renders.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	getStartDate := r.Form.Get("start_date")
	getEndDate := r.Form.Get("end_date")

	layout := "2006-01-02"
	start_date, err := time.Parse(layout, getStartDate)
	if err != nil {
		helpers.ServerError(w, err)
	}

	end_date, err := time.Parse(layout, getEndDate)
	if err != nil {
		helpers.ServerError(w, err)
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: start_date,
		EndDate:   end_date,
		RoomID:    roomID,
	}

	w.Write([]byte(fmt.Sprintf("first_name: %s, last_name: %s, email: %s",
		reservation.FirstName, reservation.LastName, reservation.Email)))

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 4, r)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		renders.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	_, err = m.DB.InsertReservation(reservation)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//restriction := models.RoomRestrictions{
	//	StartDate:     start_date,
	//	EndDate:       end_date,
	//	RoomID:        roomID,
	//	ReservationID: newReservationID,
	//	RestrictionID: 1,
	//}
	//
	//err = m.DB.InsertRoomRestriction(restriction)
	//
	//if err != nil {
	//	helpers.ServerError(w, err)
	//	return
	//}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}
