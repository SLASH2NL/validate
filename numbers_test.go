package validate_test

import (
	"testing"

	"github.com/SLASH2NL/validate"
	"github.com/stretchr/testify/require"
)

func TestNumberMin(t *testing.T) {
	err := validate.NumberMin(5)(3)
	require.NotNil(t, err)

	err = validate.NumberMin(5)(10)
	require.Nil(t, err)
}

func TestNumberMax(t *testing.T) {
	err := validate.NumberMax(5)(10)
	require.NotNil(t, err)

	err = validate.NumberMax(5)(3)
	require.Nil(t, err)
}
