package validate

import "github.com/SLASH2NL/validate/vcodes"

// Required will validate that the value is not the zero value for the type.
func Required[T comparable](value T) error {
	var x T // Create the nullable value for the type

	if value == x {
		return NewError(vcodes.Required, nil)
	}

	return nil
}

// Equal will validate that the value is equal to the expected value.
// This will not do a deep comparison.
func Equal[T comparable](expected T) Validator[T] {
	return func(value T) error {
		if value != expected {
			return NewError(vcodes.Equal, map[string]any{
				"expected": expected,
			})
		}

		return nil
	}
}

// OneOf will validate that the value is one of the accepted values.
func OneOf[T comparable](accepted ...T) Validator[T] {
	return func(value T) error {
		for _, a := range accepted {
			if value == a {
				return nil
			}
		}

		return NewError(vcodes.OneOf, map[string]any{
			"accepted": accepted,
		})
	}
}
