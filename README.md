<div align="center">
  <img src="https://github.com/scriptnull/jsonseal/assets/4211715/2bcc42dc-89b2-4844-ad29-e83682dff629" height="128px" style="max-width: 100%;" />
  <br><br>
  <span><b>jsonseal</b></span>
  <br><br>
  <span>A JSON validator for Go { ❓ 🧐 ❓ }</span>
  <br><br>

  [![Tests](https://github.com/scriptnull/jsonseal/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/scriptnull/jsonseal/actions/workflows/test.yml)
  [![Go Reference](https://pkg.go.dev/badge/github.com/scriptnull/jsonseal.svg)](https://pkg.go.dev/github.com/scriptnull/jsonseal)

  🚧 Work In Progress 🚧

</div>

&nbsp;

## Goals

- Validation errors should be human-friendly.
- Writing custom validators is a breeze. (just write a `func() error`)
- An [errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup) style API for expressing validation logic.
- A drop-in replacement for `json.Unmarshal`. (if you wish)

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
    "currency": "USD",
    "payment_mode": "card"
  }
}
```

Validation logic for the json could written as shown below:

```go
func (r *PaymentRequest) Validate() error {
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
			return fmt.Errorf("unsupported payment mode: %s", r.Payment.Mode)
		}

		return nil
	})

	return payment.Validate()
}
```

Now use `jsonseal.Unmarshal`instead of `json.Unmarshal` to inflate your struct and perform validation rules.

```go
var paymentRequest PaymentRequest

err := jsonseal.Unmarshal(paymentRequestWithInsufficientFunds, &paymentRequest)
if err != nil {
  // report error
}
```

## API

## Check Groups

Check groups are a way to group multiple checks and perform validation for them at once.

```go
var grp1 jsonseal.CheckGroup
grp1.Check(func() error { /* check condition 1 */ })
grp1.Check(func() error { /* check condition 2 */ })
err1 := grp1.Validate()

var grp2 jsonseal.CheckGroup
grp2.Check(func() error { /* check condition 1 */ })
grp2.Check(func() error { /* check condition 2 */ })
err2 := grp1.Validate()
```

## Errors
jsonseal comes with built-in error formatters for convenience.
```go
err := jsonseal.Unmarshal(paymentRequestWithInsufficientFunds, &paymentRequest)
if err != nil {
	fmt.Println("Plain error")
	fmt.Print(err)
	fmt.Println()

	fmt.Println("JSON error")
	fmt.Println(jsonseal.JSONFormat(err))
	fmt.Println()

	fmt.Println("JSON error with indent")
	fmt.Println(jsonseal.JSONIndentFormat(err, "", "  "))
	fmt.Println()
	return
}
```

But if you wish to get a Go struct that denotes all the validation errors, you could get it like this:
```go
err := jsonseal.Unmarshal(paymentRequestWithInsufficientFunds, &paymentRequest)
if err != nil {
		if validationErrors, ok := err.(*jsonseal.Errors); ok {
			fmt.Println(validationErrors)
		}
}
```
