package validate

import (
	"fmt"
)

// Validator represents a validator that can be used to validate a value.
// If a validator fails it should return an error with NewError.
// If the validator does not fail it should return nil.
type Validator[T any] func(value T) error

// Validate will run the validators on the value and return the errors grouped by the fieldName.
func Validate[F ~string, T any](
	fieldName F,
	value T,
	validators ...Validator[T],
) error {
	errs := newFieldErrors(string(fieldName))

	for _, validator := range validators {
		errs.append(validator(value))
	}

	if len(errs.errors) == 0 {
		return nil
	}

	return errs
}

// PlainValidate will run the validators on the value and return the errors.
// Callers are responsible for grouping the errors.
func PlainValidate[T any](
	value T,
	validators ...Validator[T],
) error {
	errs := newErrorList()

	for _, validator := range validators {
		errs.append(validator(value))
	}

	if len(errs.errors) == 0 {
		return nil
	}

	return errs
}

// Join will Join the errors into a single slice.
// It wil only Join errors that are of the type *validationError, *listErrors, *groupErrors.
func Join(errs ...error) error {
	verrs := newErrorList()

	for _, e := range errs {
		if e == nil {
			continue
		}

		verrs.append(e)
	}

	if len(verrs.errors) == 0 {
		return nil
	}

	return verrs
}

// If will run the validator only if the shouldRun is true.
func If[T any](shouldRun bool, validators ...Validator[T]) Validator[T] {
	if !shouldRun {
		return func(value T) error {
			return nil
		}
	}

	return func(value T) error {
		return PlainValidate(value, validators...)
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

// Slice will run the validator on each element in the slice.
func Slice[T any](value []T, validators ...Validator[T]) error {
	errs := newErrorList()

	for i, v := range value {
		errs.append(Group(fmt.Sprintf("[%d]", i), PlainValidate(v, validators...)))
	}

	if len(errs.errors) == 0 {
		return nil
	}

	return errs
}

// Map will run the validator on each value in the map and return the errors grouped by the key.
func Map[K comparable, V any](value map[K]V, validators ...Validator[V]) error {
	errs := newErrorList()

	for k, v := range value {
		errs.append(Group(fmt.Sprintf("%v", k), PlainValidate(v, validators...)))
	}

	if len(errs.errors) == 0 {
		return nil
	}

	return errs
}

// Key will run the validator on each key in the map and return the errors grouped by the key.
func Key[K comparable, V any](value map[K]V, validators ...Validator[K]) error {
	errs := newErrorList()

	for k := range value {
		errs.append(Group(fmt.Sprintf("%v", k), PlainValidate(k, validators...)))
	}

	if len(errs.errors) == 0 {
		return nil
	}

	return errs
}

// Group will add err to a a grouped error.
func Group[F ~string](field F, err error) error {
	if err == nil {
		return nil
	}

	group := newFieldErrors(string(field))
	group.append(err)

	if len(group.errors) == 0 {
		return nil
	}

	return group
}
