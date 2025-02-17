package validate

import (
	"github.com/almerlucke/go-iban/iban"
)

func IBAN[T ~string](value T) error {
	_, err := iban.NewIBAN(string(value))
	if err != nil {
		return NewError("validate.iban", nil)
	}

	return nil
}
