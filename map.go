package validate

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
			ExactPath:  v.name + "." + field,
			Violations: []Violation{{Code: CodeUnknownField}},
			Args:       Args{"key": key},
		}
	}

	violations := validate(value, validators...)
	if violations == nil {
		return nil
	}

	return Error{
		Path:       v.name + "." + field,
		ExactPath:  v.name + "." + field,
		Violations: violations,
		Args:       Args{"key": key},
	}
}

// Keys runs the validators on all keys.
func (v MapValidator[K, V]) Keys(field string, validators ...Validator[K]) error {
	var errs Errors

	for key := range v.value {
		violations := validate(key, validators...)
		if violations == nil {
			continue
		}

		errs = append(errs, Error{
			Path:       v.name + "." + field,
			ExactPath:  v.name + "." + field,
			Violations: violations,
			Args:       Args{"key": key},
		})
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Values runs the validators on all values.
func (v MapValidator[K, V]) Values(field string, validators ...Validator[V]) error {
	var errs Errors

	for key, value := range v.value {
		violations := validate(value, validators...)
		if violations == nil {
			continue
		}

		errs = append(errs, Error{
			Path:       v.name + "." + field,
			ExactPath:  v.name + "." + field,
			Violations: violations,
			Args:       Args{"key": key},
		})
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}
