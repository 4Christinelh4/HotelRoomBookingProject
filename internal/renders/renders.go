package renders

import (
	"bytes"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"my/gomodule/internal/config"
	"my/gomodule/internal/models"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) {
	td.CSRFToken = nosurf.Token(r)
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template from AppConfig
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get requested template from cache
	log.Println("tmpl: ", tmpl)

	t, ok := tc[tmpl]

	if !ok {
		log.Fatal("Cannot get template from template cache")
	}

	buf := new(bytes.Buffer)

	AddDefaultData(td, r)
	_ = t.Execute(buf, td)

	// render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing to browser", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	// get all of the files named *.page from the directory
	pages, err := filepath.Glob("templates/*.page.html")
	if err != nil {
		log.Println("filepath.Glob: ", err)
		return myCache, err
	}

	// range through all files ending with *.page
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
		log.Println("added name:", name)
	}

	return myCache, nil
}
