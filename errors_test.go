package validate_test

import (
	"fmt"
	"testing"

	"github.com/SLASH2NL/validate"
	"github.com/stretchr/testify/require"
)

func TestUnwrapping(t *testing.T) {
	err := validate.Validate(
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
