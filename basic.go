package validate

// Required will validate that the value is not the zero value for the type.
func Required[T comparable](value T) *Violation {
	var x T // Create the nullable value for the type

	if value == x {
		return &Violation{Code: CodeRequired}
	}

	return nil
}

// Equal will validate that the value is equal to the expected value.
// This will not do a deep comparison.
func Equal[T comparable](expected T) Validator[T] {
	return func(value T) *Violation {
		if value != expected {
			return &Violation{Code: CodeEqual, Args: Args{"expected": expected}}
		}

		return nil
	}
}

// OneOf will validate that the value is one of the accepted values.
func OneOf[T comparable](accepted ...T) Validator[T] {
	return func(value T) *Violation {
		for _, a := range accepted {
			if value == a {
				return nil
			}
		}

		return &Violation{Code: CodeOneOf, Args: Args{"accepted": accepted}}
	}
}
