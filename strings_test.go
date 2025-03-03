package validate_test

import (
	"testing"

	"github.com/SLASH2NL/validate"
	"github.com/stretchr/testify/require"
)

func TestPrefix(t *testing.T) {
	err := validate.Prefix("someprefix")("test")
	require.NotNil(t, err)
	require.Equal(t, "someprefix", err.(*validate.Violation).Args["prefix"])

	err = validate.Prefix("pre")("prefix")
	require.Nil(t, err)
}

func TestSuffix(t *testing.T) {
	err := validate.Suffix("somesuffix")("test")
	require.NotNil(t, err)
	require.Equal(t, "somesuffix", err.(*validate.Violation).Args["suffix"])

	err = validate.Suffix("fix")("suffix")
	require.Nil(t, err)
}

func TestEmail(t *testing.T) {
	err := validate.Email("test")
	require.NotNil(t, err)

	err = validate.Email("wvell@example.com")
	require.Nil(t, err)
}

func TestStrMin(t *testing.T) {
	err := validate.MinString(5)("test")
	require.NotNil(t, err)

	err = validate.MinString(5)("wvell")
	require.Nil(t, err)
}

func TestStrMax(t *testing.T) {
	err := validate.MaxString(5)("wvelll")
	require.NotNil(t, err)

	err = validate.MaxString(5)("test")
	require.Nil(t, err)
}

func TestStrLowercase(t *testing.T) {
	err := validate.Lowercase("Test")
	require.NotNil(t, err)

	err = validate.Lowercase("test")
	require.Nil(t, err)
}
