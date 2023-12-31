package main

import (
	"ecommerce_server/internal/db"
	"ecommerce_server/internal/models"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
    smtp struct {
        host string
        port int 
        username string
        password string
    }
    secretkey string 
    frontend string
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
	DB       models.DBModel
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf("Starting HTTP server on port %d", app.config.port)

	return srv.ListenAndServe()
}

func main() {
	var cfg config
    dsn := os.Getenv("DSN")
    user := os.Getenv("USERNAME")
    password := os.Getenv("PASSWORD")

	flag.IntVar(&cfg.port, "port", 4200, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {dev|maintance}")
	flag.StringVar(&cfg.db.dsn, "dsn", dsn, "DSN")
	flag.StringVar(&cfg.smtp.host, "smtphost", "smtp.mailtrap.io", "smtp host")
	flag.StringVar(&cfg.smtp.username, "smtpuser", user, "smtp user")
	flag.StringVar(&cfg.smtp.password, "smtppassword", password, "smtp password")
	flag.StringVar(&cfg.secretkey, "secret", "lsdfjlk2348901234asdfj", "secret key")
    flag.StringVar(&cfg.frontend, "frontend", "http://localhost:3000", "url to frontend")
	flag.IntVar(&cfg.smtp.port, "smtpport", 587, "smtp port")

	flag.Parse()

	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := db.OpenDb(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer conn.Close()

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
		DB:       models.DBModel{DB: conn},
	}

	err = app.serve()

	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}
}
