package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"https://*", "http://*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:             300,
	}))

	router.Post("/api/payment-intent", app.GetPaymentIntent)
    router.Get("/api/meeting/{id}", app.GetMeetingByID)
    router.Post("/api/create-customer-and-subscribe-to-plan", app.CreateCustomerAndSubscribeToPlan)

    router.Post("/api/authenticate", app.CreateAuthToken)
    router.Post("/api/is-authenticated", app.CheckAuthentication)
    router.Post("/api/forgot-password", app.SendPasswordResetEmail)
    router.Post("/api/reset-password", app.ResetPassword)

    router.Route("/api/admin", func(router chi.Router) {
        router.Use(app.Auth)
        
        router.Post("/virtual-terminal-succeeded", app.VirtualTerminalPaymentSucceeded)
    })

	return router
}
