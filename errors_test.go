package validate_test

import (
	"fmt"
	"testing"

	"github.com/SLASH2NL/validate"
	"github.com/stretchr/testify/require"
)

func TestIsValidationError(t *testing.T) {
	err := fmt.Errorf("some error")
	require.False(t, validate.IsValidationError(err))

	err = validate.Error{}
	require.True(t, validate.IsValidationError(err))

	err = validate.Errors{}
	require.True(t, validate.IsValidationError(err))
}

func TestLastSegment(t *testing.T) {
	err := validate.Field(
		"address.name",
		"",
		failValidatorWithCode[string]("fail"),
	)

	errs := validate.Collect(err)
	require.Equal(t, 1, len(errs))
	require.Equal(t, "address.name", errs[0].Path)
	require.Equal(t, "name", validate.LastPathSegment(errs[0].Path))

}
