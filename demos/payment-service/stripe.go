package main

import (
    "fmt"
    "errors"
)

type PaymentIntent struct {
    ID       string
    Amount   int64
    Currency string
    Status   string
}

func CreatePaymentIntent(amount int64, currency string) (*PaymentIntent, error) {
    if amount <= 0 {
        return nil, errors.New("amount must be positive")
    }
    return &PaymentIntent{
        ID:       "pi_simulated_" + fmt.Sprintf("%d", amount),
        Amount:   amount,
        Currency: currency,
        Status:   "requires_payment_method",
    }, nil
}

func ConfirmPayment(intentID string, paymentMethodID string) error {
    if intentID == "" || paymentMethodID == "" {
        return errors.New("missing required fields")
    }
    // In production: call stripe.PaymentIntents.Confirm
    return nil
}
