package handlers

import (
	"encoding/json"
	"log"
	"my/gomodule/internal/config"
	"my/gomodule/internal/driver"
	"my/gomodule/internal/forms"
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

// Reservation handles reservation form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostReservation handles reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
	}

	form := forms.New(r.PostForm)

	form.Has("first_name", r)
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		renders.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
	}

}
