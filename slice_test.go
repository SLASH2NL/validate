package validate_test

import (
	"testing"

	"github.com/SLASH2NL/validate"
	"github.com/stretchr/testify/require"
)

type testSlice struct {
	Name   string
	Amount int
}

func TestSlice(t *testing.T) {
	data := []testSlice{
		{Name: "John Deer", Amount: 9},
		{Name: "Deer John", Amount: 1},
	}

	err := validate.Slice("data", data).Items("total", func(value testSlice) *validate.Violation {
		if value.Amount < 5 {
			return &validate.Violation{Code: "max"}
		}

		return nil
	})
	require.Error(t, err)

	errs := validate.Collect(err)
	require.Equal(t, 1, len(errs))
	require.Equal(t, "data.*.total", errs[0].Path)
	require.Equal(t, "data.1.total", errs[0].ExactPath)
	require.Equal(t, "max", errs[0].Violations[0].Code)
}

func TestSliceResolve(t *testing.T) {
	data := []testSlice{
		{Name: "John Deer", Amount: 9},
		{Name: "Deer John", Amount: 1},
	}

	resolver := func(t testSlice) int { return t.Amount }

	validator := func(value int) *validate.Violation {
		if value < 5 {
			return &validate.Violation{Code: "max"}
		}

		return nil
	}

	err := validate.Slice("data", data).Items("amount", validate.Resolve(resolver, validator)...)
	require.Error(t, err)

	errs := validate.Collect(err)
	require.Equal(t, 1, len(errs))
	require.Equal(t, "data.*.amount", errs[0].Path)
	require.Equal(t, "data.1.amount", errs[0].ExactPath)
	require.Equal(t, "max", errs[0].Violations[0].Code)
}
