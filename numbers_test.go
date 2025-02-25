package validate_test

import (
	"testing"

	"github.com/SLASH2NL/validate"
	"github.com/stretchr/testify/require"
)

func TestNumberMin(t *testing.T) {
	err := validate.MinNumber(5)(3)
	require.NotNil(t, err)

	err = validate.MinNumber(5)(10)
	require.Nil(t, err)
}

func TestNumberMax(t *testing.T) {
	err := validate.MaxNumber(5)(10)
	require.NotNil(t, err)

	err = validate.MaxNumber(5)(3)
	require.Nil(t, err)
}
