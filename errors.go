package validate

import (
	"fmt"
	"strings"
)

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

// Merge merges a and b into a new Args.
// If a key exists in both a and b, the value from b is used.
func Merge(a Args, b Args) Args {
	dst := make(Args)

	for key, value := range a {
		dst[key] = value
	}

	for key, value := range b {
		dst[key] = value
	}

	return dst
}

// PrefixExactPath will prefix the exact path in the given error.
// This function accepts Error and Errors.
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
// This function accepts Error and Errors.
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
// This function accepts Error and Errors.
func PrefixBothPaths(prefix string, err error) error {
	err = PrefixPath(prefix, err)
	err = PrefixExactPath(prefix, err)

	return err
}

// LastSegment will change the Path in the given err to only use the last segment.
// This function accepts Error and Errors.
func LastSegment(err error) error {
	switch err := err.(type) {
	case Error:
		err.Path = lastSegment(err.Path)
		return err
	case Errors:
		for j, e := range err {
			e.Path = lastSegment(e.Path)
			err[j] = e
		}
		return err
	}

	return err
}

func lastSegment(s string) string {
	if i := strings.LastIndex(s, "."); i != -1 {
		return s[i+1:]
	}
	return s
}
