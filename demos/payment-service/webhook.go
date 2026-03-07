package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type StripeEvent struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
    var event StripeEvent
    if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
        http.Error(w, "bad request", 400)
        return
    }

    switch event.Type {
    case "payment_intent.succeeded":
        fmt.Println("Payment succeeded")
    case "payment_intent.payment_failed":
        fmt.Println("Payment failed")
    default:
        fmt.Printf("Unhandled event: %s\n", event.Type)
    }

    w.WriteHeader(200)
}
