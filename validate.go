package validate

// Validator represents a validator that can be used to validate a value.
// If a validator fails it should return an error with NewError.
// If the validator does not fail it should return nil.
type Validator[T any] func(value T) *Violation

// Field will run the validators on the value and return the errors grouped by the field.
func Field[T any](fieldName string, value T, validators ...Validator[T]) error {
	violations := validate(value, validators...)
	if violations == nil {
		return nil
	}

	return Error{
		Path:       fieldName,
		ExactPath:  fieldName,
		Violations: violations,
	}
}

// Join the errors into a single slice.
// It wil only Join errors that are of the type Error or Errors.
func Join(errs ...error) error {
	verrs := Errors{}

	for _, e := range errs {
		if e == nil {
			continue
		}

		switch e := e.(type) {
		case Errors:
			verrs = append(verrs, e...)
		case Error:
			verrs = append(verrs, e)
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
		return []Error{}
	}

}

// If will run the validators only if the shouldRun is true.
func If[T any](shouldRun bool, validators ...Validator[T]) []Validator[T] {
	if !shouldRun {
		return nil
	}

	return validators
}

// FailFirst will run the validators in order and return the first error.
func FailFirst[T any](validators ...Validator[T]) Validator[T] {
	return func(value T) *Violation {
		for _, validator := range validators {
			err := validator(value)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// OverridePath will override the path (not the exact path) in the given error.
// This will only override the path if the error is of the type Error or Errors.
func OverridePath(path string, err error) error {
	switch err := err.(type) {
	case Error:
		err.Path = path
		return err
	case Errors:
		for j, e := range err {
			e.Path = path
			err[j] = e
		}
		return err
	}

	return err
}

// OverrideExactPath will override the exact path in the given error.
// This will only override the exact path if the error is of the type Error or Errors.
func OverrideExactPath(path string, err error) error {
	switch err := err.(type) {
	case Error:
		err.ExactPath = path
		return err
	case Errors:
		for j, e := range err {
			e.ExactPath = path
			err[j] = e
		}
		return err
	}

	return err
}

// PrefixExactPath will prefix the exact path in the given error.
func PrefixExactPath(prefix string, err error) error {
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

// PrefixPath will prefix the path in the given error.
func PrefixPath(prefix string, err error) error {
	switch err := err.(type) {
	case Error:
		err.Path = prefix + "." + err.Path
		return err
	case Errors:
		for j, e := range err {
			e.Path = prefix + "." + e.Path
			err[j] = e
		}
		return err
	}

	return err
}

// PrefixBothPaths will prefix both the path and the exact path in the given error.
func PrefixBothPaths(prefix string, err error) error {
	err = PrefixPath(prefix, err)
	err = PrefixExactPath(prefix, err)

	return err
}

// Resolve will resolve the value and run the validators on the resolved value while preserving the original validator target.
// This is useful for validating slices or maps where you want to validate a field inside the value.
func Resolve[Original any, Resolved any](resolveFunc func(Original) Resolved, validators ...Validator[Resolved]) []Validator[Original] {
	wrapped := make([]Validator[Original], len(validators))
	for i, validator := range validators {
		validator := validator
		wrapped[i] = func(input Original) *Violation {
			resolved := resolveFunc(input)
			return validator(resolved)
		}
	}

	return wrapped
}

// validate will run the validators on the value and return the violations.
func validate[T any](
	value T,
	validators ...Validator[T],
) []Violation {
	var violations []Violation

	for _, validator := range validators {
		err := validator(value)
		if err == nil {
			continue
		}

		violations = append(violations, *err)
	}

	if len(violations) == 0 {
		return nil
	}

	return violations
}
