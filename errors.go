package validate

import (
	"errors"
	"fmt"
	"maps"
	"strings"
)

func IsValidationError(err error) bool {
	switch err.(type) {
	case Error, Errors:
		return true
	}

	// The error could be wrapped, try to unwrap it.
	var errs Errors
	if errors.As(err, &errs) {
		return true
	}

	var single Error
	return errors.As(err, &single)
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

type Violations []Violation

func (v Violations) Error() string {
	var errs []string

	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return fmt.Sprintf("violations: %v", errs)
}

type Args map[string]any

func (e Args) Add(key string, value any) Args {
	if e == nil {
		e = make(Args)
	}

	e[key] = value
	return e
}

// Merge merges a and b into a new Args.
// If a key exists in both a and b, the value from b is used.
func Merge(a Args, b Args) Args {
	if a == nil && b == nil {
		return nil
	}

	if a == nil {
		return b
	}

	if b == nil {
		return a
	}

	dst := make(Args)
	maps.Copy(dst, a)
	maps.Copy(dst, b)

	return dst
}

// LastPathSegment will return the last segment of the given path.
// It assumes the path is separated by dots.
func LastPathSegment(s string) string {
	if i := strings.LastIndex(s, "."); i != -1 {
		return s[i+1:]
	}
	return s
}

// mapError will run the mapFunc on the given err if the error is of type Error or Errors.
// This is used in the Map and Slice functions to add the correct path and args to an error if the
// validators returned a normal error instead of a list of violations.
func mapError(err error, mapFunc func(Error) Error) error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case Errors:
		for i := range e {
			e[i] = mapFunc(e[i])
		}
		return e
	case Error:
		return mapFunc(e)
	default:
		return err
	}
}

func prefixPath(path string, prefix string) string {
	if path == "" {
		return prefix
	}

	return prefix + "." + path
}
