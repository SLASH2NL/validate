package validate

import (
	"net/mail"
	"unicode"
)

func Email(value string) error {
	_, merr := mail.ParseAddress(string(value))
	if merr != nil {
		return NewError("email", nil)
	}

	return nil
}

func StrMin(length int) Validator[string] {
	return func(value string) error {
		if len(value) < length {
			return NewError("str_min", map[string]any{
				"min": length,
			})
		}

		return nil
	}
}

func StrMax(length int) Validator[string] {
	return func(value string) error {
		if len(value) > length {
			return NewError("str_max", map[string]any{
				"max": length,
			})
		}

		return nil
	}
}

func StrLowercase(value string) error {
	for _, r := range value {
		if unicode.IsUpper(r) {
			return NewError("str_lowercase", map[string]any{})
		}
	}

	return nil
}

func StrUppercase(value string) error {
	for _, r := range value {
		if unicode.IsLower(r) {
			return NewError("str_uppercase", map[string]any{})
		}
	}

	return nil
}
