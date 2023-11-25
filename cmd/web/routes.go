package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)


func (app *application) routes() http.Handler {
    router := chi.NewRouter()

    router.Get("/virtual-terminal", app.VirtualTerminal)
    router.Post("/payment-succeeded", app.PaymentSucceeded)

    return router
}
