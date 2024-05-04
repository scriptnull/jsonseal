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
// TODO
```

There are many aspects about this JSON that we would like to validate.

```go
// TODO
```