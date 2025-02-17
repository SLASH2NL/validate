package validate

import (
	"fmt"
	"iter"
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
		if !yieldErrors(err, "", "", yield) {
			return
		}
	}
}

// NewError creates a new validation error with the given code and arguments.
// It is used by validators to create a new error.
func NewError(code string, args map[string]any) error {
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
	Code  string
	Args  map[string]any
}

// FullPath returns the path and the field.
func (e Error) FullPath() string {
	return computePath(e.Path, e.Field)
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (field: %s, args: %v)", e.Code, e.Field, e.Args)
}

func yieldErrors(err error, parentPath string, field string, yield func(Error) bool) bool {
	if err == nil {
		return true
	}

	switch e := err.(type) {
	case validationError:
		return yield(Error{
			Field: field,
			Path:  parentPath,
			Code:  e.code,
			Args:  e.args,
		})
	case *errorList:
		newPath := computePath(parentPath, field)
		for _, err := range e.errors {
			if !yieldErrors(err, newPath, "", yield) {
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

type errorList struct {
	errors []error
}

func (e *errorList) append(err error) {
	if err == nil || !canAddError(err) {
		return
	}

	e.errors = append(e.errors, err)
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
	if err == nil || !canAddError(err) {
		return
	}

	e.errors = append(e.errors, err)
}

func (e *fieldErrors) Error() string {
	return fmt.Sprintf("group error for field %s: %v", e.field, e.errors)
}

type validationError struct {
	code string
	args map[string]any
}

func (e validationError) Error() string {
	return fmt.Sprintf("validation: %s (args: %v)", e.code, e.args)
}

func newErrorList() *errorList {
	return &errorList{
		errors: make([]error, 0),
	}
}

func computePath(parent string, field string) string {
	if parent == "" {
		return field
	}

	return parent + "." + field
}

func canAddError(err error) bool {
	switch err.(type) {
	case validationError, *errorList, *fieldErrors:
		return true
	}

	return true
}
