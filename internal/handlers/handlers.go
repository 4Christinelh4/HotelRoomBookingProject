package handlers

import (
	"log"
	"my/gomodule/internal/config"
	"my/gomodule/internal/driver"
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

const session_key = "remote_ip"

// var stringMap = make(map[string]string)

// Home: home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), session_key, remoteIp)

	renders.RenderTemplate(w, "home.page.html", &models.TemplateData{
		//StringMap: stringMap,
	})
}

// About
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//s := addValues(5, 8)
	//_, _ = fmt.Fprintf(w, "This is the about page, the sum is %d\n", s)
	stringMap := make(map[string]string)
	stringMap["test"] = "hello, world"

	remoteIp := m.App.Session.GetString(r.Context(), session_key)
	log.Println("remoteIP is", remoteIp)

	stringMap[session_key] = remoteIp

	renders.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
