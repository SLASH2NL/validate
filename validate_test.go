package validate_test

import (
	"testing"

	"github.com/SLASH2NL/validate"
	"github.com/SLASH2NL/validate/vcodes"
	"github.com/stretchr/testify/require"
)

func BenchmarkValidate(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := validate.Field(
			"first_name",
			"John",
			successValidator,
			validate.Equal("Peter"),
		)
		if err == nil {
			b.Fail()
		}
	}
}

func TestValidate(t *testing.T) {
	err := validate.Field(
		"first_name",
		"John",
		successValidator,
		validate.Equal("Peter"),
	)
	require.NotNil(t, err)

	errs := validate.Collect(err)

	require.Equal(t, 1, len(errs))
	require.Equal(t, "first_name", errs[0].Field)
	require.Equal(t, vcodes.Equal, errs[0].Code)
	require.Equal(t, "Peter", errs[0].Args["expected"])
}

type testSlice struct {
	Name   string
	Amount int
}

func TestSlice(t *testing.T) {
	validators := func(value testSlice) error {
		return validate.Join(
			validate.Field("name", value.Name, func(x string) error {
				if len(x) < 5 {
					return validate.NewError("min_len", nil)
				}

				return nil
			}),
			validate.Field("amount", value.Amount, func(x int) error {
				if x >= 9 {
					return validate.NewError(vcodes.NumberMax, nil)
				}

				return nil
			}),
		)
	}

	err := validate.Slice(
		[]testSlice{
			{Name: "John Deer", Amount: 9},
			{Name: "Joe Biden", Amount: 1},
		},
		validators,
	)
	require.NotNil(t, err)

	verrs := validate.Collect(err)
	require.Equal(t, 1, len(verrs))
	require.Equal(t, "amount", verrs[0].Field)
	require.Equal(t, vcodes.NumberMax, verrs[0].Code)
	require.Equal(t, "[0]", verrs[0].Path)
}

func TestMap(t *testing.T) {
	data := map[string]int{
		"John": 9,
		"Joe":  1,
	}

	err := validate.Map(
		data,
		validate.Key("John", validate.NumberMax(5)),
		validate.Key("Joe", validate.NumberMax(5)),
	)
	require.NotNil(t, err)
	errs := validate.Collect(err)
	require.Equal(t, 1, len(errs))
	require.Equal(t, "John", errs[0].Field)
	require.Equal(t, vcodes.NumberMax, errs[0].Code)
}

func TestIf(t *testing.T) {
	// Validate true condition should run the validator.
	err := validate.Field(
		"first_name",
		"John",
		validate.If(true, failValidator[string]),
	)
	require.NotNil(t, err)

	errs := validate.Collect(err)
	require.Equal(t, 1, len(errs))
	require.Equal(t, vcodes.Code("fail"), errs[0].Code)

	// Validate false condition should not run the validator.
	err = validate.Field[string](
		"first_name",
		"John",
		validate.If(false, failValidator[string]),
	)
	require.Nil(t, err)
}

func TestFailFirst(t *testing.T) {
	// Validate true condition should run only 1 validator.
	errs := validate.Collect(validate.Field(
		"first_name",
		"John",
		validate.FailFirst(
			failValidatorWithCode[string]("first"),
			failValidatorWithCode[string]("second"),
		),
	))
	require.NotNil(t, errs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, "first_name", errs[0].Field)
	require.Equal(t, vcodes.Code("first"), errs[0].Code)
}

func TestJoin(t *testing.T) {
	err := validate.Join(
		validate.Field(
			"name",
			"",
			failValidatorWithCode[string]("name.fail"),
		),
		validate.Group("address",
			validate.Field(
				"street",
				"",
				validate.FailFirst(
					failValidatorWithCode[string]("street.fail.first"),
					failValidatorWithCode[string]("street.fail.second"),
				),
			),
		),
	)

	errs := validate.Collect(err)

	require.Equal(t, 2, len(errs))
	require.Equal(t, "", errs[0].Path)
	require.Equal(t, "name", errs[0].Field)
	require.Equal(t, vcodes.Code("name.fail"), errs[0].Code)

	require.Equal(t, "address", errs[1].Path)
	require.Equal(t, "street", errs[1].Field)
	require.Equal(t, vcodes.Code("street.fail.first"), errs[1].Code)
}

func TestGroup(t *testing.T) {
	err := validate.Group("person",
		validate.Group(
			"address",
			validate.Field(
				"street",
				"Some street 123",
				failValidator[string],
			),
		),
	)

	errs := validate.Collect(err)
	require.Equal(t, 1, len(errs))
	require.Equal(t, "person.address", errs[0].Path)
	require.Equal(t, "street", errs[0].Field)
	require.Equal(t, vcodes.Code("fail"), errs[0].Code)
}

func failValidatorWithCode[T any](code vcodes.Code) validate.Validator[T] {
	return func(value T) error {
		return validate.NewError(code, nil)
	}
}

func failValidator[T any](value T) error {
	return validate.NewError("fail", nil)
}

func successValidator[T any](value T) error {
	return nil
}
