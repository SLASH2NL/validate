package validate_test

import (
	"testing"

	"github.com/SLASH2NL/validate"
	"github.com/stretchr/testify/require"
)

func TestIf(t *testing.T) {
	// Validate true condition should run the validator.
	err := validate.Field(
		"first_name",
		"John",
		validate.If(true, failValidator[string])...,
	)
	require.NotNil(t, err)

	errs := validate.Collect(err)
	require.Equal(t, 1, len(errs))
	require.Equal(t, "first_name", errs[0].Path)
	require.Equal(t, "fail", errs[0].Violations[0].Code)

	// Validate false condition should not run the validator.
	err = validate.Field[string](
		"first_name",
		"John",
		validate.If(false, failValidator[string])...,
	)
	require.Nil(t, err)
}

func TestFailFirst(t *testing.T) {
	err := validate.Field(
		"first_name",
		"John",
		validate.FailFirst(
			failValidatorWithCode[string]("first"),
			failValidatorWithCode[string]("second"),
		),
	)
	require.NotNil(t, err)
	errs := validate.Collect(err)
	require.Equal(t, 1, len(errs))
	require.Equal(t, "first_name", errs[0].Path)
	require.Equal(t, "first", errs[0].Violations[0].Code)
}

func TestJoin(t *testing.T) {
	err := validate.Join(
		validate.Field(
			"name",
			"",
			failValidatorWithCode[string]("name.fail"),
		),
		validate.Field(
			"iban",
			"invalid",
			failValidatorWithCode[string]("iban.fail"),
		),
	)

	errs := validate.Collect(err)

	require.Equal(t, 2, len(errs))
	require.Equal(t, "name", errs[0].Path)
	require.Equal(t, "name.fail", errs[0].Violations[0].Code)

	require.Equal(t, "iban", errs[1].Path)
	require.Equal(t, "iban.fail", errs[1].Violations[0].Code)
}

func TestOverride(t *testing.T) {
	t.Run("override", func(t *testing.T) {
		err := validate.Join(
			validate.Field(
				"name",
				"",
				failValidatorWithCode[string]("name.fail"),
				failValidatorWithCode[string]("name.fail2"),
			),
			validate.Field(
				"iban",
				"invalid",
				failValidatorWithCode[string]("iban.fail"),
			),
		)

		err = validate.Override("test", err)

		errs := validate.Collect(err)
		require.Equal(t, 2, len(errs))
		require.Equal(t, "test", errs[0].Path)
		require.Equal(t, "name.fail", errs[0].Violations[0].Code)
		require.Equal(t, "name.fail2", errs[0].Violations[1].Code)
	})

	t.Run("override nil", func(t *testing.T) {
		err := validate.Override("test", nil)
		require.Nil(t, err)
	})
}

func failValidatorWithCode[T any](code string) validate.Validator[T] {
	return func(value T) *validate.Violation {
		return &validate.Violation{Code: code}
	}
}

func failValidator[T any](value T) *validate.Violation {
	return &validate.Violation{Code: "fail"}
}
