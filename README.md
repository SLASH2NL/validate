# Golang validation library
Validate is a simple golang validation library that is focussed on simplicity and typed validators.

[![GoDev](https://pkg.go.dev/badge/github.com/SLASH2NL/validate)](https://pkg.go.dev/github.com/SLASH2NL/validate)

## Usage
See [the example](https://pkg.go.dev/github.com/SLASH2NL/validate) for the usage of the library.

## Creating custom validators
A validator is a simple function that takes a value and returns an error if the value is invalid.
The error should be created with validate.NewError(code, args).

```go
// Returns an error if the number is not 42.
func Is42(x int) error {
    if x != 42 {
        return validate.NewError("Is42", "The number must be 42.")
    }

    return nil
}
```