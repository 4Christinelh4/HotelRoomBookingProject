package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"log"
	"my/gomodule/internal/config"
	"my/gomodule/internal/driver"
	"my/gomodule/internal/handlers"
	"my/gomodule/internal/helpers"
	"my/gomodule/internal/models"
	"my/gomodule/internal/renders"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8044"

var infoLog *log.Logger
var errorLog *log.Logger

var app config.AppConfig
var session *scs.SessionManager

func run() (*driver.DB, error) {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// read flags
	inProduction := flag.Bool("production", true, "App is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	//dbUser := flag.String("dbuser", "", "Database user name")
	//dbPass := flag.String("dbpassword", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	//dbSSL := flag.String("dbssl", "disable", "Database ssl")

	flag.Parse()
	if *dbName == "" {
		fmt.Println("Missing require flags")
		os.Exit(1)
	}

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	app.InProduction = *inProduction

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

	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s",
		*dbHost, *dbPort, *dbName)

	// "host=localhost port=5432 dbname=bookings"
	db, err := driver.ConnectSQL(connectionString)
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

	app.UseCache = *useCache
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
	defer close(app.MailChan)

	listenForMail()

	fmt.Println(fmt.Sprintf("Starting application on port 8044"))
	srv := &http.Server{
		Addr:    portNumber,
		Handler: RoutesApp(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
		log.Fatal("Listen and serve error")
	}
}
