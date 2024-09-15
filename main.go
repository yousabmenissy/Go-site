package main

import (
	"flag"
	"net/http"
	"site/data"
	"site/internal"
	"site/internal/mailer"
	"text/template"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/lib/pq"
)

type Config struct {
	Domain string
	Port   string
	DB     struct {
		DbHost     string
		DbName     string
		DbUser     string
		DbPort     string
		DbPassword string
	}
	Smtp struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
}

type Application struct {
	config          *Config
	models          *data.Models
	templatesChache map[string]*template.Template
	sessionManager  *scs.SessionManager
	mailer          mailer.Mailer
	logger          *internal.ConsoleLogger
}

func main() {
	var config Config

	flag.StringVar(&config.Domain, "domain", "localhost", "the http port the server will listen to")
	flag.StringVar(&config.Port, "port", ":8080", "the http port the server will listen to")

	flag.StringVar(&config.DB.DbHost, "dbHost", "localhost", "the the database host")
	flag.StringVar(&config.DB.DbName, "DbName", "site", "the name of the database used")
	flag.StringVar(&config.DB.DbUser, "DbUser", "postgres", "the name of the database user")
	flag.StringVar(&config.DB.DbPort, "DbPort", "5432", "the database port")
	flag.StringVar(&config.DB.DbPassword, "DbPassword", "postgres", "the database user's password")

	flag.StringVar(&config.Smtp.Host, "smpt-host", "<your-smpt-host>", "SMTP host")
	flag.StringVar(&config.Smtp.Username, "smpt-username", "<your-smpt-username>", "SMTP username")
	flag.StringVar(&config.Smtp.Password, "smpt-password", "<your-smpt-password>", "SMTP password")
	flag.StringVar(&config.Smtp.Sender, "smpt-sender", "<the-sender-email>", "SMTP sender")
	flag.IntVar(&config.Smtp.Port, "smpt-port", 587, "SMPT port")

	flag.Parse()

	logger := internal.NewConsoleLogger()

	db, err := config.OpenConnection()
	if err != nil {
		logger.LogError.Println(err)
		db.Close()
	}
	defer db.Close()

	templateChache, err := NewTemplateCache()
	if err != nil {
		logger.LogError.Fatalln(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	app := &Application{
		config:          &config,
		models:          data.NewModels(db),
		templatesChache: templateChache,
		sessionManager:  sessionManager,
		mailer:          mailer.New(config.Smtp.Host, config.Smtp.Port, config.Smtp.Username, config.Smtp.Password, config.Smtp.Sender),
		logger:          logger,
	}

	server := &http.Server{
		Addr:     config.Port,
		Handler:  app.routes(),
		ErrorLog: logger.LogError,
	}

	logger.LogInfo.Println("Starting server on http://localhost" + config.Port)
	logger.LogError.Panic(server.ListenAndServe())
}
