package validate

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode"
)

func Email(value string) error {
	_, merr := mail.ParseAddress(string(value))
	if merr != nil {
		return &Violation{Code: CodeEmail}
	}

	return nil
}

func Regex(re *regexp.Regexp) Validator[string] {
	return func(value string) error {
		if !re.MatchString(value) {
			return &Violation{Code: CodeRegex, Args: Args{"pattern": re.String()}}
		}

		return nil
	}
}

func MinString(length int) Validator[string] {
	return func(value string) error {
		if len(value) < length {
			return &Violation{Code: CodeStringMin, Args: Args{"min": length}}
		}

		return nil
	}
}

func MaxString(length int) Validator[string] {
	return func(value string) error {
		if len(value) > length {
			return &Violation{Code: CodeStringMax, Args: Args{"max": length}}
		}

		return nil
	}
}

func Lowercase(value string) error {
	for _, r := range value {
		if unicode.IsUpper(r) {
			return &Violation{Code: CodeLowercase}
		}
	}

	return nil
}

func Uppercase(value string) error {
	for _, r := range value {
		if unicode.IsLower(r) {
			return &Violation{Code: CodeUppercase}
		}
	}

	return nil
}

func Prefix(prefix string) Validator[string] {
	return func(value string) error {
		if !strings.HasPrefix(value, prefix) {
			return &Violation{Code: "prefix", Args: Args{"prefix": prefix}}
		}

		return nil
	}
}

func Suffix(suffix string) Validator[string] {
	return func(value string) error {
		if !strings.HasSuffix(value, suffix) {
			return &Violation{Code: "suffix", Args: Args{"suffix": suffix}}
		}

		return nil
	}
}
