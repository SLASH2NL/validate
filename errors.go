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

type Args map[string]any

func (e Args) Merge(key string, value any) Args {
	if e == nil {
		e = make(Args)
	}

	e[key] = value
	return e
}
