package validate_test

import (
	"testing"

	"github.com/SLASH2NL/validate"
	"github.com/stretchr/testify/require"
)

func TestIBAN(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		violation := validate.IBAN("invalid")
		require.NotNil(t, violation)
		require.Equal(t, validate.CodeIBAN, violation.(*validate.Violation).Code)
	})

	t.Run("valid", func(t *testing.T) {
		violation := validate.IBAN("NL91ABNA0417164300")
		require.Nil(t, violation)
	})
}
