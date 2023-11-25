package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"
const cssVersion = "1"

type config struct {
    port int
    env string
    api string
    db struct {
        dsn string
    }
    stripe struct {
        secret string
        key string
    }
}

type application struct {
    config config
    infoLog *log.Logger
    errorLog *log.Logger
    templateCache map[string]*template.Template
    version string
}

func (app *application) serve() error {
    srv := &http.Server {
        Addr: fmt.Sprintf(":%d", app.config.port),
        Handler: app.routes(),
        IdleTimeout: 30 * time.Second,
        ReadTimeout: 10 * time.Second,
        ReadHeaderTimeout: 5 * time.Second,
        WriteTimeout: 5 * time.Second,
    }

    app.infoLog.Printf("Starting HTTP server on port %d", app.config.port)

    return srv.ListenAndServe()
}

func main() {
    var cfg config

    flag.IntVar(&cfg.port, "port", 3000, "Server port to listen on")
    flag.StringVar(&cfg.env, "env", "development", "Application environment {dev|production}")
    flag.StringVar(&cfg.api, "api", "http://localhost:4200", "URL to api")

    flag.Parse()

    cfg.stripe.key = os.Getenv("STRIPE_KEY")
    cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
        
    tc := make(map[string]*template.Template) 

    app := &application{
        config: cfg,
        infoLog: infoLog,
        errorLog: errorLog,
        templateCache: tc,
        version: version,
    }

    err := app.serve()
    
    if err != nil {
       app.errorLog.Println(err) 
       log.Fatal(err)
    }
}
