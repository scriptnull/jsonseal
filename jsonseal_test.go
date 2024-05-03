package jsonseal_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/scriptnull/jsonseal"
)

type PaymentRequest struct {
	Payments []Payment `json:"payments"`
}

func (p *PaymentRequest) Validate() error {
	var validate jsonseal.ValidateAll

	for idx, payment := range p.Payments {
		payment := payment

		validate.Fieldf("payments[%d].amount", idx).Check(func() error {
			if payment.Amount <= 0 {
				return errors.New("amount should be greater than 0")
			}
			return nil
		})

		validate.Fieldf("payments[%d].currency", idx).Check(func() error {
			if !slices.Contains(SupportedCurrencies, payment.Currency) {
				return errors.New("unsupported currency")
			}
			return nil
		})

		validate.Check(func() error {
			if !slices.Contains(SupportedPaymentModes, payment.Mode) {
				return errors.New("unspported payment mode")
			}

			return nil
		})

		validate.Check(func() error {
			if payment.Detail == nil {
				return errors.New("expected valid payment details")
			}

			if ok, err := payment.Detail.Valid(); !ok {
				return err
			}

			return nil
		})
	}

	return validate.Validate()
}

type Payment struct {
	Amount   float64       `json:"amount"`
	Currency Currency      `json:"currency"`
	Mode     PaymentMode   `json:"payment_mode"`
	Detail   PaymentDetail `json:"payment_detail"`
}

func (p *Payment) UnmarshalJSON(data []byte) error {
	// Do minimal parsing to inflate the struct
	// Avoid user input validation at this level

	type payment struct {
		Amount   float64         `json:"amount"`
		Currency Currency        `json:"currency"`
		Mode     PaymentMode     `json:"payment_mode"`
		Detail   json.RawMessage `json:"payment_detail"`
	}

	var pi payment
	err := json.Unmarshal(data, &pi)
	if err != nil {
		return err
	}

	p.Amount = pi.Amount
	p.Currency = pi.Currency
	p.Mode = pi.Mode

	switch p.Mode {
	case Card:
		p.Detail = new(CardDetail)
	case Upi:
		p.Detail = new(UPIDetail)
	}
	if p.Detail != nil {
		err = json.Unmarshal(pi.Detail, &(p.Detail))
		if err != nil {
			return err
		}
	}

	return nil
}

type Currency string

const (
	INR Currency = "INR"
	USD Currency = "USD"
)

var (
	SupportedCurrencies = []Currency{INR, USD}
)

type PaymentMode string

const (
	Card PaymentMode = "card"
	Upi  PaymentMode = "upi"
)

var (
	SupportedPaymentModes = []PaymentMode{Card, Upi}
)

type PaymentDetail interface {
	Valid() (bool, error)
}

type UPIDetail struct {
	UPIID string `json:"upi_id"`
}

func (u *UPIDetail) Valid() (bool, error) {
	upiDetails := strings.Split(u.UPIID, "@")

	if len(upiDetails) != 2 {
		return false, errors.New("expected format: <username>@<bankname> for UPI ID")
	}

	// more validation for user name and bank name

	return true, nil
}

type CardDetail struct {
	Number  string `json:"card_number"`
	ExpDate string `json:"exp_date"`
}

func (c *CardDetail) Valid() (bool, error) {
	if len(c.Number) != 16 {
		return false, errors.New("card number should have 16 numbers")
	}

	// maybe use a library like https://github.com/durango/go-credit-card that takes care of card validation

	return true, nil
}

func TestValidateAll(t *testing.T) {
	validPayments, err := os.ReadFile("testcases/valid_payments.json")
	if err != nil {
		t.Fatal(err)
	}
	invalidPaymentsWithOneError, err := os.ReadFile("testcases/invalid_payments_1.json")
	if err != nil {
		t.Fatal(err)
	}
	tt := []struct {
		name        string
		jsonContent []byte
		expectedErr error
	}{
		{
			name:        "valid json",
			jsonContent: validPayments,
			expectedErr: nil,
		},
		{
			name:        "invalid json with one error inside an array object",
			jsonContent: invalidPaymentsWithOneError,
			expectedErr: &jsonseal.Errors{
				[]jsonseal.Error{
					{
						Fields: []string{"payments[0].amount"},
						Err:    errors.New("amount should be greater than 0"),
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var paymentRequest PaymentRequest
			err = json.Unmarshal(tc.jsonContent, &paymentRequest)
			if err != nil {
				t.Fatal(err)
			}

			err = paymentRequest.Validate()
			if tc.expectedErr != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Errorf("expected %s, got %s", tc.expectedErr.Error(), err.Error())
				}

				gotJsonErr, _ := json.Marshal(err)
				expectedJsonErr, _ := json.Marshal(tc.expectedErr)
				if !bytes.Equal(gotJsonErr, expectedJsonErr) {
					t.Errorf("expected %s, got %s", string(expectedJsonErr), string(gotJsonErr))
				}
			}
		})
	}
}
