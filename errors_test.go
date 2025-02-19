package validate_test

import (
	"fmt"
	"testing"

	"github.com/SLASH2NL/validate"
	"github.com/SLASH2NL/validate/vcodes"
	"github.com/stretchr/testify/require"
)

func TestUnwrapping(t *testing.T) {
	err := validate.Field(
		"first_name",
		"John",
		failValidator,
	)

	// Wrap the error.
	err = fmt.Errorf("this error is wrapped: %w", err)

	// Collect should unwrap the error.
	errs := validate.Collect(err)
	require.Len(t, errs, 1)
	require.Equal(t, "first_name", errs[0].Field)
}

func TestCollectUnknown(t *testing.T) {
	err := validate.NewError(vcodes.Equal, nil)

	errs := validate.Collect(err)
	require.Len(t, errs, 1)
	require.Equal(t, validate.UnknownField, errs[0].Field)
	require.Equal(t, vcodes.Equal, errs[0].Code)
}
