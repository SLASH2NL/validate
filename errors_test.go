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
