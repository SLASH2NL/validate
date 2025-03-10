package validate

import (
	"fmt"
)

func Map[K comparable, V any](name string, value map[K]V) MapValidator[K, V] {
	return MapValidator[K, V]{
		name:  name,
		value: value,
	}
}

type MapValidator[K comparable, V any] struct {
	name  string
	value map[K]V
}

// Key runs the validators on the value of the key.
// If the key does not exist, it will return an unknown.field violation.
func (v MapValidator[K, V]) Key(field string, key K, validators ...Validator[V]) error {
	value, ok := v.value[key]
	if !ok {
		return Error{
			Path:       v.name + "." + field,
			ExactPath:  v.name + "." + fmt.Sprintf("%v", key) + "." + field,
			Violations: []Violation{{Code: CodeUnknownField}},
			Args:       Args{"key": key},
		}
	}

	var verrs Errors

	violations, err := validate(value, validators...)
	if err != nil {
		// It could be that the validators returned an Error or Errors. If so we map it with the correct paths.
		if isValidationError(err) {
			switch err := err.(type) {
			case Error:
				verrs = verrs.Merge(prefixMapError(err, v.name, field, key))
			case Errors:
				verrs = verrs.MergeAll(err.mapErrors(func(err Error) Error {
					return prefixMapError(err, v.name, field, key)
				}))
			}
		} else {
			return err
		}
	}

	if len(violations) > 0 {
		verrs = append(verrs, Error{
			Path:       v.name + "." + field,
			ExactPath:  v.name + "." + fmt.Sprintf("%v", key) + "." + field,
			Violations: violations,
			Args:       Args{"key": key},
		})
	}

	if len(verrs) == 0 {
		return nil
	}

	return verrs
}

// Keys runs the validators on all keys.
func (v MapValidator[K, V]) Keys(field string, validators ...Validator[K]) error {
	var verrs Errors

	for key := range v.value {
		violations, err := validate(key, validators...)
		if err != nil {
			// It could be that the validators returned an Error or Errors. If so we map it with the correct paths.
			if isValidationError(err) {
				switch err := err.(type) {
				case Error:
					verrs = verrs.Merge(prefixMapError(err, v.name, field, key))
				case Errors:
					verrs = verrs.MergeAll(err.mapErrors(func(err Error) Error {
						return prefixMapError(err, v.name, field, key)
					}))
				}
			} else {
				return err
			}
		}

		if len(violations) > 0 {
			verrs = append(verrs, Error{
				Path:       v.name + "." + field,
				ExactPath:  v.name + "." + fmt.Sprintf("%v", key) + "." + field,
				Violations: violations,
				Args:       Args{"key": key},
			})
		}
	}

	if len(verrs) == 0 {
		return nil
	}

	return verrs
}

// Values runs the validators on all values.
func (v MapValidator[K, V]) Values(field string, validators ...Validator[V]) error {
	var verrs Errors

	for key, value := range v.value {
		violations, err := validate(value, validators...)
		if err != nil {
			// It could be that the validators returned an Error or Errors. If so we map it with the correct paths.
			if isValidationError(err) {
				switch err := err.(type) {
				case Error:
					verrs = verrs.Merge(prefixMapError(err, v.name, field, key))
				case Errors:
					verrs = verrs.MergeAll(err.mapErrors(func(err Error) Error {
						return prefixMapError(err, v.name, field, key)
					}))
				}
			} else {
				return err
			}
		}

		if len(violations) > 0 {
			verrs = append(verrs, Error{
				Path:       v.name + "." + field,
				ExactPath:  v.name + "." + fmt.Sprintf("%v", key) + "." + field,
				Violations: violations,
				Args:       Args{"key": key},
			})
		}
	}

	if len(verrs) == 0 {
		return nil
	}

	return verrs
}

func prefixMapError(err Error, name string, field string, key any) Error {
	err.Path = prefixPath(err.Path, name+"."+field)
	err.ExactPath = prefixPath(err.ExactPath, name+"."+fmt.Sprintf("%v", key)+"."+field)
	err.Args = err.Args.Add("key", key)
	return err
}
