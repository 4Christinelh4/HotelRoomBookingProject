package handlers

import (
	"encoding/json"
	"log"
	"my/gomodule/internal/config"
	"my/gomodule/internal/driver"
	"my/gomodule/internal/helpers"
	"my/gomodule/internal/models"
	"my/gomodule/internal/renders"
	"my/gomodule/internal/repository"
	"my/gomodule/internal/repository/dbrepo"
	"net/http"
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
	renders.Template(w, r, "home.page.html", &models.TemplateData{})
}

// About
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "about.page.html", &models.TemplateData{})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "generals.page.html", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "majors.html", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "search-availability.page.html", &models.TemplateData{})
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
	renders.Template(w, r, "contact.page.html", &models.TemplateData{})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//sd := r.Form.Get("start_date")
	//ed := r.Form.Get("end_date")
	//
	//reservation := {
	//
	//}
}
