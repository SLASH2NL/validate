package validate

import (
	"github.com/SLASH2NL/validate/vcodes"
	"golang.org/x/exp/constraints"
)

func NumberMin[T constraints.Integer | constraints.Float](min T) Validator[T] {
	return func(value T) error {
		if value < min {
			return NewError(vcodes.NumberMin, map[string]any{
				"min": min,
			})
		}

		return nil
	}
}

func NumberMax[T constraints.Integer | constraints.Float](max T) Validator[T] {
	return func(value T) error {
		if value > max {
			return NewError(vcodes.NumberMax, map[string]any{
				"max": max,
			})
		}

		return nil
	}
}
