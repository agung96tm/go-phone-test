package main

import (
	"database/sql"
	"flag"
	"github.com/agung96tm/go-phone-test/internal/models"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	addr string
	DB   struct {
		dsn string
	}
	cors struct {
		trustedOrigins []string
	}
}

func DefaultConfig() Config {
	return Config{
		cors: struct{ trustedOrigins []string }{
			trustedOrigins: []string{
				"http://localhost:3000",
				"http://localhost:4000",
				"http://localhost:5000",
			},
		},
	}
}

type application struct {
	models   *models.Models
	infoLog  *log.Logger
	errorLog *log.Logger
	config   Config
}

func main() {
	cfg := DefaultConfig()

	flag.StringVar(&cfg.addr, "addr", ":8000", "HTTP network address")
	flag.StringVar(&cfg.DB.dsn, "db-dsn", "postgres://phone_user:phone_password@localhost:5432/phone_db?sslmode=disable", "Database DSN")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := initDB(cfg.DB.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := application{
		models:   models.New(db),
		infoLog:  infoLog,
		errorLog: errorLog,
		config:   cfg,
	}

	srv := &http.Server{
		Addr:         cfg.addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on: %s\n", cfg.addr)
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
