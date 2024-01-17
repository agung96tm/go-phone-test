package main

import (
	"database/sql"
	"flag"
	"github.com/agung96tm/go-phone-test/internal/authentication"
	"github.com/agung96tm/go-phone-test/internal/models"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	models        *models.Models
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	googleOauth2  *authentication.GoogleOauth2
}

func main() {
	cfg := DefaultConfig()

	flag.StringVar(&cfg.SecretKey, "secret-key", "foobar", "")
	flag.StringVar(&cfg.Addr, "addr", ":8000", "HTTP network Address")
	flag.StringVar(&cfg.DB.dsn, "db-dsn", "postgres://phone_user:phone_password@localhost:5432/phone_db?sslmode=disable", "Database DSN")

	flag.StringVar(&cfg.googleOauth2.ClientID, "oauth2-google-clientid", "", "Google Oauth2 Client ID")
	flag.StringVar(&cfg.googleOauth2.ClientSecret, "oauth2-google-clientsecret", "", "Google Oauth2 Client Secret")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := initDB(cfg.DB.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := application{
		models:        models.New(db),
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: templateCache,
		googleOauth2: authentication.NewGoogleOauth2(
			cfg.googleOauth2.RedirectURL,
			cfg.googleOauth2.SendTokenUrl,
			cfg.googleOauth2.ClientID,
			cfg.googleOauth2.ClientSecret,
		),
	}

	srv := &http.Server{
		Addr:         cfg.Addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on: %s\n", cfg.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func initDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
