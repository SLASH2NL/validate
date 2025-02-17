# Golang validation library
Validate is a simple golang validation library that is focussed on simplicity and typed validators.

## Usage
Validation can be nested by using Join or Group to combine multiple validators.
```go
// Validate a complex struct.
type SomeStruct struct {
    FirstName string
    LastName string
    Address struct {
        Street string
    }
}

func (x SomeStruct) validate() error {
    return validate.Join(
        validate.Validate(
            "first_name",
            x.FirstName,
            validate.Required,
        ),
        validate.Validate(
            "last_name",
            x.LastName,
            validate.FailFirst(
                validate.Required,
                validate.StrMax(255),
            ),
        ),
        validate.Group("address", validate.Validate(
            "street",
            x.Address.Street,
            validate.Required,
        )),
    )
}
```

Errors can be extracted from the result of the validation to perform custom error handling and/or translation.
```go
errs := validate.Collect(err)

// errs will contain a list of Error.
type Error struct {
	Field string
	Path  string
	Code  string
	Args  map[string]any // Optional arguments from the validator.
}

// For the above example this can contain(if all validators fail):
// {Field: "first_name", Path: "", Code: "required", Args: nil}
// {Field: "last_name", Path: "", Code: "required", Args: nil}
// {Field: "street", Path: "address", Code: "required", Args: nil}
```

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