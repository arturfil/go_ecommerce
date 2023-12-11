package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)


func (app *application) routes() http.Handler {
    router := chi.NewRouter()
    router.Use(SessionLoad)

    router.Get("/", app.VirtualTerminal)
    router.Post("/virtual-terminal-payment-succeeded", app.VirtualPaymentSucceeded)
    router.Get("/virtual-terminal-receipt", app.VirtualTerminalReceipt)
    router.Get("/session/{id}", app.ChargeOnce)

    router.Post("/payment-succeeded", app.PaymentSucceeded)
    router.Get("/receipt", app.Receipt)

    router.Get("/plans/bronze", app.BronzePlan)

    fileServer := http.FileServer(http.Dir("./static")) 

    router.Handle("/static/*", http.StripPrefix("/static", fileServer))

    return router
}
