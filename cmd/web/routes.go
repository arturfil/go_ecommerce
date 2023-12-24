package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)


func (app *application) routes() http.Handler {
    router := chi.NewRouter()
    router.Use(SessionLoad)

    router.Route("/admin", func(router chi.Router) {
        router.Use(app.Auth)
        router.Get("/", app.VirtualTerminal)
    })

    // router.Post("/virtual-terminal-payment-succeeded", app.VirtualPaymentSucceeded)
    // router.Get("/virtual-terminal-receipt", app.VirtualTerminalReceipt)

    router.Get("/", app.HomePage)

    router.Get("/meeting/{id}", app.ChargeOnce)

    router.Post("/payment-succeeded", app.PaymentSucceeded)
    router.Get("/receipt", app.Receipt)

    router.Get("/plans/bronze", app.BronzePlan)
    router.Get("/receipt/bronze", app.BronzePlanReceipt)

    router.Get("/login", app.LoginPage)
    router.Post("/login", app.PostLoginPage)
    router.Get("/logout", app.Logout)
    router.Get("/forgot-password", app.ForgotPassword)
    router.Get("/reset-password", app.ShowResetPassword)


    fileServer := http.FileServer(http.Dir("./static")) 

    router.Handle("/static/*", http.StripPrefix("/static", fileServer))

    return router
}
