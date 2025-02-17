package validate

import (
	"golang.org/x/exp/constraints"
)

func NumberMin[T constraints.Integer | constraints.Float](min T) Validator[T] {
	return func(value T) error {
		if value < min {
			return NewError("min", map[string]any{
				"min": min,
			})
		}

		return nil
	}
}

func NumberMax[T constraints.Integer | constraints.Float](max T) Validator[T] {
	return func(value T) error {
		if value > max {
			return NewError("max", map[string]any{
				"max": max,
			})
		}

		return nil
	}
}
