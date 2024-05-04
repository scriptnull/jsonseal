package jsonseal_test

import (
	"errors"
	"fmt"
	"slices"

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
		"payment_mode": "card"
	}
}
`)

type SimplePaymentRequest struct {
	AccountID string   `json:"account_id"`
	Balance   float64  `json:"balance"`
	Currency  Currency `json:"currency"`
	Payment   struct {
		Amount   float64     `json:"amount"`
		Currency Currency    `json:"currency"`
		Mode     PaymentMode `json:"mode"`
	} `json:"payment"`
}

func (r *SimplePaymentRequest) Validate() error {
	var payment jsonseal.CheckGroup

	payment.Check(func() error {
		if r.Payment.Currency != r.Currency {
			return errors.New("payment not allowed to different currency")
		}

		if r.Payment.Amount > r.Balance {
			return errors.New("insufficent balance")
		}

		return nil
	})

	payment.Check(func() error {
		if !slices.Contains(SupportedPaymentModes, r.Payment.Mode) {
			return fmt.Errorf("unsupported payment mode")
		}

		return nil
	})

	return payment.Validate()
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
