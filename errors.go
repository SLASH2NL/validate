package validate

import (
	"errors"
	"fmt"
	"iter"

	"github.com/SLASH2NL/validate/vcodes"
)

const (
	UnknownField = "unknown"
)

// Collect returns all the internal errors as a slice of Error.
func Collect(err error) Errors {
	var errors Errors

	for e := range CollectIter(err) {
		errors = append(errors, e)
	}

	return errors
}

// CollectIter will iterate over all the internal errors and convert them to Error and return them.
// If an error can not be converted to an Error, it will be ignored.
func CollectIter(err error) iter.Seq[Error] {
	return func(yield func(Error) bool) {
		// If the error is wrapped we try to unwrap it.
		switch err.(type) {
		case validationError, *errorList, *fieldErrors:
		default:
			// Try to find errors in order.
			var list *errorList
			var field *fieldErrors
			var verr validationError
			if errors.As(err, &list) {
				err = list
			} else if errors.As(err, &field) {
				err = field
			} else if errors.As(err, &verr) {
				err = verr
			}
		}

		if !yieldErrors(err, "", "", yield) {
			return
		}
	}
}

// NewError creates a new validation error with the given code and arguments.
// It is used by validators to create a new error.
func NewError(code vcodes.Code, args ErrArgs) error {
	return validationError{
		code: code,
		args: args,
	}
}

// A simple collection type to hold all the errors.
// It is a placeholder to possibly add more functionality in the future.
type Errors []Error

// Error is the error type returned by the Collect method.
// It is the only exposed error type.
type Error struct {
	Field string
	Path  string
	Code  vcodes.Code
	Args  ErrArgs
}

type ErrArgs map[string]any

func (e ErrArgs) Merge(key string, value any) ErrArgs {
	if e == nil {
		e = make(ErrArgs)
	}

	e[key] = value
	return e
}

// FullPath returns the path and the field.
func (e Error) FullPath() string {
	return computePath(e.Path, e.Field)
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (path: %s, field: %s, args: %v)", e.Code, e.Path, e.Field, e.Args)
}

func yieldErrors(err error, parentPath string, field string, yield func(Error) bool) bool {
	if err == nil {
		return true
	}

	switch e := err.(type) {
	case validationError:
		if field == "" {
			field = UnknownField
		}

		return yield(Error{
			Field: field,
			Path:  parentPath,
			Code:  e.code,
			Args:  e.args,
		})
	case *errorList:
		for _, err := range e.errors {
			if !yieldErrors(err, parentPath, field, yield) {
				return false
			}
		}
	case *fieldErrors:
		newPath := parentPath
		if field != "" {
			newPath = computePath(parentPath, field)
		}

		for _, err := range e.errors {
			if !yieldErrors(err, newPath, e.field, yield) {
				return false
			}
		}
	}

	return true
}

func newErrorList() *errorList {
	return &errorList{
		errors: make([]error, 0),
	}
}

type errorList struct {
	errors []error
}

func (e *errorList) append(err error) {
	if err == nil {
		return
	}

	switch err := err.(type) {
	case *errorList:
		e.errors = append(e.errors, err.errors...)
	case *fieldErrors:
		// Only add field errors if there are any errors.
		if len(err.errors) > 0 {
			e.errors = append(e.errors, err)
		}
	case validationError:
		e.errors = append(e.errors, err)
	default:
		e.errors = append(e.errors, NewError(vcodes.Unknown, map[string]any{
			"err": err,
		}))
	}
}

func (e *errorList) Error() string {
	return fmt.Sprintf("validation errors: %v", e.errors)
}

func newFieldErrors(field string) *fieldErrors {
	return &fieldErrors{
		field:  field,
		errors: make([]error, 0),
	}
}

type fieldErrors struct {
	field  string
	errors []error `` // Can be slice of GroupError again.
}

func (e *fieldErrors) append(err error) {
	if err == nil {
		return
	}

	switch err := err.(type) {
	case *errorList:
		e.errors = append(e.errors, err.errors...)
	case *fieldErrors:
		// Only add field errors if there are any errors.
		if len(err.errors) > 0 {
			e.errors = append(e.errors, err)
		}
	case validationError:
		e.errors = append(e.errors, err)
	default:
		e.errors = append(e.errors, NewError(vcodes.Unknown, map[string]any{
			"err": err,
		}))
	}
}

func (e *fieldErrors) Error() string {
	return fmt.Sprintf("group error for field %s: %v", e.field, e.errors)
}

type validationError struct {
	code vcodes.Code
	args map[string]any
}

func (e validationError) Error() string {
	return fmt.Sprintf("validation: %s (args: %v)", e.code, e.args)
}

func computePath(parent string, field string) string {
	if parent == "" {
		return field
	}

	return parent + "." + field
}
