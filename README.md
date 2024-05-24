<div align="center">
  <img src="https://github.com/scriptnull/jsonseal/assets/4211715/59caaf7a-2def-4556-8370-4213a4358d87" height="128px" style="max-width: 100%;" />
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
			return errors.New("insufficient balance")
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
- [Check Groups](#check-groups)
- [Errors](#errors)
- [Fields](#fields)
- [Drop-in Replacements](#drop-in-replacements)
- [Unknown Field Suggestions](#unknown-field-suggestions)

### Check Groups

Check groups are a way to group multiple checks and perform validation for them at once.

```go
var grp1 jsonseal.CheckGroup
grp1.Check(func() error { /* check condition 1 */ })
grp1.Check(func() error { /* check condition 2 */ })
err1 := grp1.Validate()

var grp2 jsonseal.CheckGroup
grp2.Check(func() error { /* check condition 1 */ })
grp2.Check(func() error { /* check condition 2 */ })
err2 := grp2.Validate()
```

### Errors

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

An example error message that is returned by `jsonseal.JSONIndentFormat` looks like

```js
{
  "errors": [
    {
      "error": "insufficient balance"
    },
    {
      "error": "unsupported payment mode: neft"
    }
  ]
}
```

### Fields

JSON fields could be associated with the validation errors like this:

```go
payment.Field("payment.mode").Check(func() error {
	if !slices.Contains(SupportedPaymentModes, r.Payment.Mode) {
		return fmt.Errorf("unsupported payment mode: %s", r.Payment.Mode)
	}

	return nil
})
```

The above code associates the json field `payment.mode` with any error that arises from the `Check` block attached to it.

```js
// before calling Field()
{
  "error": "unsupported payment mode: neft"
}

// after calling Field()
{
  "fields": [
    "payment.mode"
  ],
  "error": "unsupported payment mode: neft"
}
```

A method called `Fieldf` is available to help with cases like `payments.Fieldf("payments[%d].amount", idx)` (while trying to associate array element as a field).

An error could be asscoiated with multiple different fields by chaining `Field` or `Fieldf`.

```go
users.Field("sender.id").Field("receiver.id").Check(AreFriends())
```

This will make sure to associate both fields with the error in case of the validation error.

```js
{
  "fields": [
    "sender.id",
    "receiver.id"
  ],
  "error": "sender and receiver are not friends"
}
```

### Drop-in Replacements

jsonseal provides drop-in replacements for a few things in [encoding/json](https://pkg.go.dev/encoding/json) package. This is to ensure API compatibility and seamless migration experience.

- `jsonseal.Unmarshal` could be used in the place of `json.Unmarshal`
- `jsonseal.Decoder` could be used in the place of `json.Decoder`
  ```go
  err = jsonseal.NewDecoder(data).Decode(&v)
  ```

If you wish to ensure that `jsonseal.Validator` interface is satisfied by validation input, you could use the below alternatives:
- `jsonseal.UnmarshalValidate` could be used instead of `jsonseal.Unmarshal`.
- `jsonseal.DecodeValidate` could be used instead of `jsonseal.Decode`.

Alternatively, you could also do the following to ensure the compile time guarantee.

```go
var _ jsonseal.Validator = &PaymentRequest{}
```

### Unknown Field Suggestions

It might be useful to validate if the JSON data contains only the fields that are expected by the struct to which it is decoded to.

Example: A user sends `{"expires": 50}` as the JSON data but our code expects it to be `{"expires_in": 50}`. If you are using `json` package, you might enable this validation by calling `DisallowUnknownFields()` on the `json.Decoder`. That will give you an error like `json: unknown field "expires"`.

jsonseal provides `WithUnknownFieldSuggestion()` method which takes the error message to the next level by suggesting the right field name based on the [Levenshtein Distance](https://en.wikipedia.org/wiki/Levenshtein_distance) between the wrongly typed field name and all possible field names of the struct that we are decoding to.

```go
type Data struct {
  ExpiresIn      int    `json:"expires_in"`
  Balance        int    `json:"balance,omitempty"`
  PrivateField   string `json:"-"`
}
var d Data
err := jsonseal.NewDecoder(data).WithUnknownFieldSuggestion().Decode(&d)
if err != nil {
  fmt.Println(jsonseal.JSONIndentFormat(err, "", "  "))
}
```

This gives the following error

```json
{
  "errors": [
    {
      "fields": ["expires"],
      "error": "unknown field. Did you mean \"expires_in\""
    }
  ]
}
```

Thanks to [this blog post](https://brandur.org/disallow-unknown-fields) for providing inspiration and motivation for this feature in jsonseal :pray:.
