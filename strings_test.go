package validate_test

import (
	"testing"

	"github.com/SLASH2NL/validate"
	"github.com/stretchr/testify/require"
)

func TestEmail(t *testing.T) {
	err := validate.Email("test")
	require.NotNil(t, err)

	err = validate.Email("wvell@example.com")
	require.Nil(t, err)
}

func TestStrMin(t *testing.T) {
	err := validate.StrMin(5)("test")
	require.NotNil(t, err)

	err = validate.StrMin(5)("wvell")
	require.Nil(t, err)
}

func TestStrMax(t *testing.T) {
	err := validate.StrMax(5)("wvelll")
	require.NotNil(t, err)

	err = validate.StrMax(5)("test")
	require.Nil(t, err)
}

func TestStrLowercase(t *testing.T) {
	err := validate.StrLowercase("Test")
	require.NotNil(t, err)

	err = validate.StrLowercase("test")
	require.Nil(t, err)
}
