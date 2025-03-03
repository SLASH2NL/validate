package validate_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/SLASH2NL/validate"
	"github.com/stretchr/testify/require"
)

func TestReturnError(t *testing.T) {
	exception := errors.New("some exception")

	err := validate.Join(
		validate.Field(
			"first_name",
			"John",
			failValidator[string],
		),
		validate.Field(
			"iban",
			"invalid",
			func(value string) error {
				return exception
			},
		),
	)
	require.NotNil(t, err)
	require.Equal(t, exception, err)
}

func TestIf(t *testing.T) {
	// Validate true condition should run the validators.
	err := validate.Field(
		"first_name",
		"John",
		validate.If(
			true,
			failValidatorWithCode[string]("first"),
			failValidatorWithCode[string]("second"),
			successValidator[string],
		),
	)
	require.NotNil(t, err)

	errs := validate.Collect(err)
	require.Equal(t, 1, len(errs))
	require.Equal(t, 2, len(errs[0].Violations))
	require.Equal(t, "first_name", errs[0].Path)
	require.Equal(t, "first", errs[0].Violations[0].Code)
	require.Equal(t, "second", errs[0].Violations[1].Code)

	// Validate false condition should not run the validator.
	err = validate.Field[string](
		"first_name",
		"John",
		validate.If(false, failValidator[string]),
	)
	require.Nil(t, err)
}

func TestAnd(t *testing.T) {
	t.Run("all validators pass", func(t *testing.T) {
		err := validate.Field(
			"first_name",
			"John",
			validate.And(
				successValidator[string],
				successValidator[string],
			),
		)
		require.Nil(t, err)
	})

	t.Run("a validator fails", func(t *testing.T) {
		err := validate.Field(
			"first_name",
			"John",
			validate.And(
				failValidatorWithCode[string]("failing"),
				successValidator[string],
			),
		)
		require.Nil(t, err)
	})

	t.Run("all validators fail", func(t *testing.T) {
		err := validate.Field(
			"first_name",
			"John",
			validate.And(
				failValidatorWithCode[string]("first"),
				failValidatorWithCode[string]("second"),
			),
		)
		require.Error(t, err)
		errs := validate.Collect(err)
		require.Equal(t, 1, len(errs))
		require.Equal(t, "first_name", errs[0].Path)
		require.Equal(t, "first", errs[0].Violations[0].Code)
		require.Equal(t, "second", errs[0].Violations[1].Code)
	})
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

func TestJoinMergesErrors(t *testing.T) {
	err := validate.Join(
		validate.Field(
			"name",
			"",
			failValidatorWithCode[string]("name.fail"),
		),
		validate.Field(
			"name",
			"",
			failValidatorWithCode[string]("name.second"),
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
	require.Equal(t, "name.second", errs[0].Violations[1].Code)

	require.Equal(t, "iban", errs[1].Path)
	require.Equal(t, "iban.fail", errs[1].Violations[0].Code)
}

func TestCollect(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
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
	})

	t.Run("unwrap error", func(t *testing.T) {
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

		err = fmt.Errorf("some error: %w", err)

		errs := validate.Collect(err)
		require.Equal(t, 2, len(errs))
	})
}

func TestReplaceIfErr(t *testing.T) {
	err := validate.Field(
		"name",
		"",
		failValidatorWithCode[string]("name.fail"),
		failValidatorWithCode[string]("name.fail2"),
	)

	override := errors.New("some error")

	err = validate.ReplaceIfErr(err, override)
	require.ErrorIs(t, err, override)
}

func TestGroup(t *testing.T) {
	err := validate.Field(
		"name",
		"",
		failValidatorWithCode[string]("name.fail"),
		failValidatorWithCode[string]("name.fail2"),
	)

	err = validate.Group("prefix", err)

	errs := validate.Collect(err)
	require.Equal(t, 1, len(errs))
	require.Equal(t, "name", errs[0].Path)
	require.Equal(t, "prefix.name", errs[0].ExactPath)
}

func failValidatorWithCode[T any](code string) validate.Validator[T] {
	return func(value T) error {
		return &validate.Violation{Code: code}
	}
}

func failValidator[T any](value T) error {
	return &validate.Violation{Code: "fail"}
}

func successValidator[T any](value T) error {
	return nil
}
