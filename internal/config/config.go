package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
	"my/gomodule/internal/models"
)

// AppConfig holds the application config
type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
