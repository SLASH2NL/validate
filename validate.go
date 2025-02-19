package validate

import (
	"fmt"
)

// Validator represents a validator that can be used to validate a value.
// If a validator fails it should return an error with NewError.
// If the validator does not fail it should return nil.
type Validator[T any] func(value T) error

// Validate will run the validators on the value and return the errors.
// Errors are not scoped to a field, callers are responsible for grouping the errors or using Field.
func Validate[T any](
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

// Field will run the validators on the value and return the errors grouped by the field.
func Field[F ~string, T any](fieldName F, value T, validators ...Validator[T]) error {
	err := Validate(value, validators...)
	if err == nil {
		return nil
	}

	return Group(fieldName, err)
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
		return Validate(value, validators...)
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
		errs.append(Group(fmt.Sprintf("[%d]", i), Validate(v, validators...)))
	}

	if len(errs.errors) == 0 {
		return nil
	}

	return errs
}

// Map will run the every KeyValidator and return the errors grouped by the key defined in the KeyValidator.
func Map[K comparable, V any](value map[K]V, mapValidators ...MapValidator[K, V]) error {
	errs := newErrorList()

	for _, v := range mapValidators {
		value, ok := value[v.matchKey]
		if !ok {
			errs.append(Group(fmt.Sprintf("%v", v.matchKey), NewError(UnknownField, nil)))
			continue
		}

		errs.append(Group(fmt.Sprintf("%v", v.matchKey), Validate(value, v.validators...)))
	}

	if len(errs.errors) == 0 {
		return nil
	}

	return errs
}

// MapKeys will run the validators on every key.
func MapKeys[K comparable, V any](data map[K]V, validators ...Validator[K]) error {
	errs := newErrorList()

	for key := range data {
		errs.append(Group(fmt.Sprintf("%v", key), Validate(key, validators...)))
	}

	if len(errs.errors) == 0 {
		return nil
	}

	return errs
}

// MapValues will run the validators on every value.
func MapValues[K comparable, V any](data map[K]V, validators ...Validator[V]) error {
	errs := newErrorList()

	for key, value := range data {
		errs.append(Group(fmt.Sprintf("%v", key), Validate(value, validators...)))
	}

	if len(errs.errors) == 0 {
		return nil
	}

	return errs
}

// Key will return a MapValidator that runs the validator on the value of the given key.
func Key[K comparable, V any](key K, validators ...Validator[V]) MapValidator[K, V] {
	return MapValidator[K, V]{
		matchKey:   key,
		validators: validators,
	}
}

// Group will group the error by field.
// If the given error is a validationError it will be transformed to:
//
//	fieldErrors{ field: "myField", errors: []error{err} }
//
// Using it this way will make Collect return the error like this:
//
//	[]Error{ { Field: "myField", Path: "", Code: the-code, Args: the-args } }
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

// GroupValidators will group the validators by field.
// This can be used for adding a field when validators are run in Slice or Map.
// Example:
//
//	err := validate.Slice(
//		[]string{"John", "Doe"},
//		validate.GroupValidators("first_name", validate.Required),
//	)
//
// This will return an error with the Path: "[0]" and the Field: "first_name".
func GroupValidators[F ~string, T any](field F, value T, validators ...Validator[T]) Validator[T] {
	return func(value T) error {
		return Group(field, Validate(value, validators...))
	}
}

type MapValidator[K comparable, V any] struct {
	matchKey   K
	validators []Validator[V]
}
