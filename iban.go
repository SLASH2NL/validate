package validate

import (
	"github.com/SLASH2NL/validate/vcodes"
	"github.com/almerlucke/go-iban/iban"
)

func IBAN[T ~string](value T) error {
	_, err := iban.NewIBAN(string(value))
	if err != nil {
		return NewError(vcodes.IBAN, nil)
	}

	return nil
}
