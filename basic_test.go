package validate_test

import (
	"testing"

	"github.com/SLASH2NL/validate"
	"github.com/stretchr/testify/require"
)

func TestRequired(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		err := validate.Required(0)
		require.NotNil(t, err)

		err = validate.Required(1)
		require.Nil(t, err)
	})

	t.Run("string", func(t *testing.T) {
		err := validate.Required("")
		require.NotNil(t, err)

		err = validate.Required("got value")
		require.Nil(t, err)
	})

	t.Run("float", func(t *testing.T) {
		err := validate.Required(0.0)
		require.NotNil(t, err)

		err = validate.Required(0.5)
		require.Nil(t, err)
	})

	t.Run("bool", func(t *testing.T) {
		err := validate.Required(false)
		require.NotNil(t, err)

		err = validate.Required(true)
		require.Nil(t, err)
	})
}

func TestOneOf(t *testing.T) {
	err := validate.OneOf(2, 3, 4)(1)
	require.NotNil(t, err)

	err = validate.OneOf(1, 2, 3, 4)(1)
	require.Nil(t, err)
}

func TestNot(t *testing.T) {
	err := validate.Not(1)(1)
	require.NotNil(t, err)

	err = validate.Not(2)(1)
	require.Nil(t, err)
}

func TestEqual(t *testing.T) {
	err := validate.Equal(1)(1)
	require.Nil(t, err)

	err = validate.Equal(2)(1)
	require.NotNil(t, err)
}
