<div align="center">
  <img src="https://github.com/scriptnull/jsonseal/assets/4211715/2bcc42dc-89b2-4844-ad29-e83682dff629" height="128px" style="max-width: 100%;" />
  <br><br>
  <span><b>jsonseal</b></span>
  <br><br>
  <span>A JSON validator for Go { ❓ 🧐 ❓ }</span>
  <br><br>

  [![Tests](https://github.com/scriptnull/jsonseal/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/scriptnull/jsonseal/actions/workflows/test.yml)

</div>

&nbsp;

## Design choices

- Writing custom validators should be easy: just write a `func() error`, you got one!
- In case of a validation error, provide a way to capture information about where the error occured.

## Installation

```
go get github.com/scriptnull/jsonseal
```

## Example

Consider the following JSON, that could arrive in a web request for performing payments.

```js
{
  "account_id": "3ee7b5eb-f3fc-4f0b-9e01-8d7a0fa76f0b",
  "balance": 15,
  "currency": "USD",
  "payment": {
		"amount": 50,
		"currency": "USD"
	}
}
```

A Go struct could be defined like below to parse JSON and perform some validations on top of it.

```go
type PaymentRequest struct {
	AccountID string   `json:"account_id"`
	Balance   float64  `json:"balance"`
	Currency  Currency `json:"currency"`
	Payment   struct {
		Amount   float64  `json:"amount"`
		Currency Currency `json:"currency"`
	} `json:"payment"`
}

func (r *PaymentRequest) Validate() error {
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
```

Now use `jsonseal.Unmarshal`instead of `json.Unmarshal` to inflate your struct and perform validation rules.

```go
err := jsonseal.Unmarshal(paymentRequestWithInsufficientFunds, &paymentRequest)
if err != nil {
  // report error
}
```