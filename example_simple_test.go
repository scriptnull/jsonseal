package jsonseal_test

import (
	"errors"
	"fmt"

	"github.com/scriptnull/jsonseal"
)

var paymentRequestWithInsufficientFunds = []byte(`
{
  "account_id": "3ee7b5eb-f3fc-4f0b-9e01-8d7a0fa76f0b",
  "balance": 15,
  "currency": "USD",
  "payment": {
		"amount": 50,
		"currency": "USD"
	}
}
`)

type SimplePaymentRequest struct {
	AccountID string   `json:"account_id"`
	Balance   float64  `json:"balance"`
	Currency  Currency `json:"currency"`
	Payment   struct {
		Amount   float64  `json:"amount"`
		Currency Currency `json:"currency"`
	} `json:"payment"`
}

func (r *SimplePaymentRequest) Validate() error {
	var payments jsonseal.CheckGroup

	payments.Check(func() error {
		if r.Payment.Currency == r.Currency && r.Payment.Amount > r.Balance {
			return errors.New("insufficent balance")
		}

		return nil
	})

	payments.Check(func() error {
		if r.Payment.Currency != r.Currency {
			return errors.New("payment not allowed to different currency")
		}

		return nil
	})

	return payments.Validate()
}

func Example_simple() {
	var paymentRequest SimplePaymentRequest

	err := jsonseal.Unmarshal(paymentRequestWithInsufficientFunds, &paymentRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Output:
	// insufficent balance
}
