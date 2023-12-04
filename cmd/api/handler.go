package main

import (
	"ecommerce_server/internal/cards"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type stripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id"`
}

func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorLog.Println(err)
		return
	}


	amount, err := strconv.Atoi(payload.Amount)
    log.Println("AMOUNT ->", amount)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: payload.Currency,
	}

	okay := true

	pIntent, msg, err := card.Charge(payload.Currency, amount)
	if err != nil {
		okay = false
	}

	if okay {

		out, err := json.MarshalIndent(pIntent, "", " ")
		if err != nil {
			app.errorLog.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)

	} else {

        j := jsonResponse{
            OK: false,
            Message: msg,
            Content: "",
        }

        out, err := json.MarshalIndent(j, "", "   ")
        if err != nil {
            app.errorLog.Println(err)
        }

        w.Header().Set("Content-Type", "application/json")
        w.Write(out)

	}

}

func (app *application) GetSessionByID(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    sessionID, _ := strconv.Atoi(id)

    session, err := app.DB.GetSession(sessionID)
    if err != nil {
        app.errorLog.Println(err)
        return 
    }

    out, err := json.MarshalIndent(session, "", " ")
    if err != nil {
        app.errorLog.Println(err)
        return 
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(out)
}
