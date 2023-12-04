package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)


func (app *application) routes() http.Handler {
    router := chi.NewRouter()
    router.Use(SessionLoad)

    router.Get("/", app.VirtualTerminal)
    router.Post("/payment-succeeded", app.PaymentSucceeded)
    router.Get("/session/{id}", app.ChargeOnce)

    fileServer := http.FileServer(http.Dir("./static")) 

    router.Handle("/static/*", http.StripPrefix("/static", fileServer))

    return router
}
