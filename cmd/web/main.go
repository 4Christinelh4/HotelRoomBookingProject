package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"log"
	"my/gomodule/internal/config"
	"my/gomodule/internal/driver"
	"my/gomodule/internal/handlers"
	"my/gomodule/internal/helpers"
	"my/gomodule/internal/renders"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8080"

var infoLog *log.Logger
var errorLog *log.Logger

var app config.AppConfig
var session *scs.SessionManager

func run() (*driver.DB, error) {
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to DB")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings")
	if err != nil {
		log.Fatal("Cannot connect to DB bookings")
		return nil, err
	}
	log.Println("Connected to DB")

	tc, err := renders.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}

	log.Println("Created template cache")

	app.UseCache = true
	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)

	renders.NewTemplates(&app)
	handlers.NewHandlers(repo)
	helpers.NewHelper(&app)

	return db, nil
}

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Starting application on port 8080"))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: RoutesApp(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Listen and serve error")
	}
}
