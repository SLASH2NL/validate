package validate

import "errors"

// Validator represents a validator that can be used to validate a value.
// If a validator fails it should return an new Violation.
// If there is an unexpected exception a normal error should be returned. This error
// will bubble up and be returned to the caller.
type Validator[T any] func(value T) error

// Field will run the validators on the value and return the errors grouped by the field.
// If a violation returned a non Violation that is returned as exception error.
func Field[T any](fieldName string, value T, validators ...Validator[T]) error {
	violations, err := validate(value, validators...)
	if err != nil {
		return err
	}

	if violations == nil {
		return nil
	}

	return Error{
		Path:       fieldName,
		ExactPath:  fieldName,
		Violations: violations,
	}
}

// Join the errors into a single slice and merge all errors with the same exact path.
// It wil only Join errors that are of the type Error or Errors.
func Join(errs ...error) error {
	verrs := Errors{}

	for _, e := range errs {
		if e == nil {
			continue
		}

		switch e := e.(type) {
		case Errors:
			for _, err := range e {
				verrs = verrs.Merge(err)
			}
		case Error:
			verrs = verrs.Merge(e)
		default:
			// If we encountered an exception we just return that.
			return e
		}
	}

	if len(verrs) == 0 {
		return nil
	}

	return verrs
}

// Collect will collect the Errors from the given error.
func Collect(err error) []Error {
	switch e := err.(type) {
	case Errors:
		return e
	case Error:
		return []Error{e}
	default:
		// The error could be wrapped, try to unwrap it.
		var errs Errors
		if errors.As(err, &errs) {
			return errs
		}

		var single Error
		if errors.As(err, &single) {
			return []Error{single}
		}

		return []Error{}
	}

}

// And will run all validators and only return an error if all validators error.
func And[T any](validators ...Validator[T]) Validator[T] {
	return func(value T) error {
		var allViolations Violations

		for _, validator := range validators {
			violations, err := validate(value, validator)
			if err != nil {
				// If an a non violation error is returned we will bubble it up.
				return err
			}

			// We can return early if there are no violations in one validator.
			if len(violations) == 0 {
				return nil
			}

			allViolations = append(allViolations, violations...)
		}

		if len(allViolations) == 0 {
			return nil
		}

		return allViolations
	}
}

// If will run the validators only if the shouldRun is true.
func If[T any](shouldRun bool, validators ...Validator[T]) Validator[T] {
	if !shouldRun {
		return func(value T) error {
			return nil
		}
	}

	return func(v T) error {
		violations, err := validate(v, validators...)
		if err != nil {
			return err
		}

		if violations == nil {
			return nil
		}

		return Violations(violations)
	}
}

// FailFirst will run the validators in order and return the first error.
func FailFirst[T any](validators ...Validator[T]) Validator[T] {
	return func(value T) error {
		for _, validator := range validators {
			err := validator(value)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// Resolve will resolve the value and run the validators on the resolved value while preserving the original validator target.
// This is useful for validating slices or maps where you want to validate a field inside the value.
func Resolve[Original any, Resolved any](resolveFunc func(Original) Resolved, validators ...Validator[Resolved]) []Validator[Original] {
	wrapped := make([]Validator[Original], len(validators))
	for i, validator := range validators {
		validator := validator
		wrapped[i] = func(input Original) error {
			resolved := resolveFunc(input)
			return validator(resolved)
		}
	}

	return wrapped
}

// Group will prefix the exact path in the given error.
// This function accepts Error and Errors.
func Group(prefix string, err error) error {
	switch err := err.(type) {
	case Error:
		err.ExactPath = prefix + "." + err.ExactPath
		return err
	case Errors:
		for j, e := range err {
			e.ExactPath = prefix + "." + e.ExactPath
			err[j] = e
		}
		return err
	}

	return err
}

// validate will run the validators on the value and return the violations.
// If a validator returns a non *Violation error it will return that error and discard the violations.
func validate[T any](
	value T,
	validators ...Validator[T],
) ([]Violation, error) {
	var violations []Violation

	for _, validator := range validators {
		err := validator(value)
		if err == nil {
			continue
		}

		switch err := err.(type) {
		case Violations:
			violations = append(violations, err...)
		case *Violation:
			violations = append(violations, *err)
		default:
			return nil, err
		}
	}

	if len(violations) == 0 {
		return nil, nil
	}

	return violations, nil
}
