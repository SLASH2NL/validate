package validate

import "fmt"

// Slice will run the validators on each element in the slice.
func Slice[F ~string, T any](name F, value []T) SliceValidator[T] {
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
	var verrs Errors

	path := fmt.Sprintf("%s.*", v.name)
	for i, value := range v.value {
		violations, err := validate(value, validators...)
		if err != nil {
			// It could be that the validators returned an Error or Errors. If so we map it with the correct paths.
			if isValidationError(err) {
				switch err := err.(type) {
				case Error:
					verrs = verrs.Merge(prefixSliceError(err, v.name, field, i))
				case Errors:
					verrs = verrs.MergeAll(err.mapErrors(func(err Error) Error {
						return prefixSliceError(err, v.name, field, i)
					}))
				}
			} else {
				return err
			}
		}

		if len(violations) > 0 {
			verrs = append(verrs, Error{
				Path:       path + "." + field,
				ExactPath:  fmt.Sprintf("%s.%d.%s", v.name, i, field),
				Violations: violations,
				Args:       Args{"index": i},
			})
		}
	}

	if len(verrs) == 0 {
		return nil
	}

	return verrs
}

func prefixSliceError(err Error, name string, field string, index int) Error {
	err.Path = prefixPath(err.Path, fmt.Sprintf("%s.*.%s", name, field))
	err.ExactPath = prefixPath(err.ExactPath, fmt.Sprintf("%s.%d.%s", name, index, field))
	err.Args = err.Args.Add("index", index)
	return err
}
