package validate

import "reflect"

// NotNil will return an error if value is nil.
func NotNil[T any](value T) error {
	// Do a reflect check to see if the value is nil.
	vof := reflect.ValueOf(value)
	switch vof.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		if vof.IsNil() {
			return &Violation{Code: CodeNotNil}
		}
	}

	return nil
}

// Not will validate that the value is not the given value.
func Not[T comparable](not T) Validator[T] {
	return func(value T) error {
		if value == not {
			return &Violation{Code: CodeNot, Args: Args{"not": not}}
		}

		return nil
	}
}

// Required will validate that the value is not the zero value for the type.
func Required[T comparable](value T) error {
	var x T // Create the nullable value for the type

	if value == x {
		return &Violation{Code: CodeRequired}
	}

	return nil
}

// Equal will validate that the value is equal to the expected value.
// This will not do a deep comparison.
func Equal[T comparable](expected T) Validator[T] {
	return func(value T) error {
		if value != expected {
			return &Violation{Code: CodeEqual, Args: Args{"expected": expected}}
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

		return &Violation{Code: CodeOneOf, Args: Args{"accepted": accepted}}
	}
}
