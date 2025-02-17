package validate

import (
	"net/mail"
	"regexp"
	"unicode"
)

func Email(value string) error {
	_, merr := mail.ParseAddress(string(value))
	if merr != nil {
		return NewError("validate.string.email", nil)
	}

	return nil
}

func StrRegexRaw(pattern string, re *regexp.Regexp) Validator[string] {
	return func(value string) error {
		if !re.MatchString(value) {
			return NewError("validate.string.regex", map[string]any{
				"pattern": pattern,
			})
		}

		return nil
	}
}

func StrRegex(pattern string) Validator[string] {
	re, err := regexp.Compile(pattern)

	return func(value string) error {
		if err != nil {
			return NewError("validate.string.regex.invalid", map[string]any{
				"pattern": pattern,
				"error":   err.Error(),
			})
		}

		if !re.MatchString(value) {
			return NewError("validate.string.regex", map[string]any{
				"pattern": pattern,
			})
		}

		return nil
	}
}

func StrMin(length int) Validator[string] {
	return func(value string) error {
		if len(value) < length {
			return NewError("validate.string.min", map[string]any{
				"min": length,
			})
		}

		return nil
	}
}

func StrMax(length int) Validator[string] {
	return func(value string) error {
		if len(value) > length {
			return NewError("validate.string.max", map[string]any{
				"max": length,
			})
		}

		return nil
	}
}

func StrLowercase(value string) error {
	for _, r := range value {
		if unicode.IsUpper(r) {
			return NewError("validate.string.lowercase", map[string]any{})
		}
	}

	return nil
}

func StrUppercase(value string) error {
	for _, r := range value {
		if unicode.IsLower(r) {
			return NewError("validate.string.uppercase", map[string]any{})
		}
	}

	return nil
}
