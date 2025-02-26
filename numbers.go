package validate

import (
	"golang.org/x/exp/constraints"
)

func MinNumber[T constraints.Integer | constraints.Float](min T) Validator[T] {
	return func(value T) error {
		if value < min {
			return &Violation{Code: CodeNumberMin, Args: Args{"min": min}}
		}

		return nil
	}
}

func MaxNumber[T constraints.Integer | constraints.Float](max T) Validator[T] {
	return func(value T) error {
		if value > max {
			return &Violation{Code: CodeNumberMax, Args: Args{"max": max}}
		}

		return nil
	}
}
