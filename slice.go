package validate

import "fmt"

// Slice will run the validators on each element in the slice.
func Slice[F ~string, T any](name F, value []T, validators ...Validator[T]) SliceValidator[T] {
	return SliceValidator[T]{
		name:  string(name),
		value: value,
	}
}

type SliceValidator[T any] struct {
	name  string
	value []T
}

func (v SliceValidator[T]) Items(field string, validators ...Validator[T]) error {
	var errs Errors

	path := fmt.Sprintf("%s.*", v.name)
	for i, value := range v.value {
		violations, err := validate(value, validators...)
		if err != nil {
			return err
		}

		if violations == nil {
			continue
		}

		errs = append(errs, Error{
			Path:       path + "." + field,
			ExactPath:  fmt.Sprintf("%s.%d.%s", v.name, i, field),
			Violations: violations,
			Args:       Args{"index": i},
		})
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}
