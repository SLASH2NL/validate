package validate

import "fmt"

func IsValidationError(err error) bool {
	switch err.(type) {
	case Error, Errors:
		return true
	}

	return false
}

type Errors []Error

// Merge merges the given error into the errors.
// If there is already an error with the same exact path, it will merge the violations.
// Otherwise it is added to the errors.
func (e Errors) Merge(errs Error) Errors {
	for i, err := range e {
		if err.ExactPath == errs.ExactPath {
			e[i].Violations = append(e[i].Violations, errs.Violations...)
			return e
		}
	}

	return append(e, errs)
}

func (e Errors) Error() string {
	var errs []string

	for _, err := range e {
		errs = append(errs, err.Error())
	}

	return fmt.Sprintf("validation errors: %v", errs)
}

type Error struct {
	Path       string
	ExactPath  string
	Args       Args
	Violations []Violation
}

func (e Error) Error() string {
	return fmt.Sprintf("validation error for exact path: %s, path: %s, args: %v, violations: %v", e.ExactPath, e.Path, e.Args, e.Violations)
}

type Violation struct {
	Code string
	Args Args
}

func (v Violation) Error() string {
	return fmt.Sprintf("violation code: %s, args: %v", v.Code, v.Args)
}

type Args map[string]any

func (e Args) Add(key string, value any) Args {
	if e == nil {
		e = make(Args)
	}

	e[key] = value
	return e
}

func (e Args) Merge(from Args) Args {
	if e == nil {
		if from == nil {
			return make(Args)
		}

		return from
	}

	for key, value := range from {
		e[key] = value
	}

	return e
}
