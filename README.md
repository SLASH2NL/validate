# Golang validation library
Validate is a simple golang validation library that is focussed on simplicity and typed validators.

[![GoDev](https://pkg.go.dev/badge/github.com/SLASH2NL/validate)](https://pkg.go.dev/github.com/SLASH2NL/validate)

## Usage
See `examples_test.go` for usage.

## Creating custom validators
A validator is a simple function that takes a value and returns an error if the value is invalid.
If the error is a violation it should return a `validate.Violation` error.
If the error is an exception that should be handled by the caller it should return a normal error.

```go
// Returns an error if the number is not 42.
func Is42(x int) error {
    if x != 42 {
        return &validate.Violdation{ Code: "Is42" }
    }

    return nil
}
```