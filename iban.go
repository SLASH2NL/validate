package validate

import (
	"github.com/almerlucke/go-iban/iban"
)

func IBAN(value string) *Violation {
	_, err := iban.NewIBAN(string(value))
	if err != nil {
		return &Violation{Code: CodeIBAN}
	}

	return nil
}
