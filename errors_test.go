package validate_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/SLASH2NL/validate"
	"github.com/stretchr/testify/require"
)

func TestIsValidationError(t *testing.T) {
	err := fmt.Errorf("some error")
	require.False(t, validate.IsValidationError(err))

	err = validate.Error{}
	require.True(t, validate.IsValidationError(err))

	err = validate.Errors{}
	require.True(t, validate.IsValidationError(err))

	// It could be wrapped.
	err = validate.Errors{}
	wrap := fmt.Errorf("wrapped: %w", err)
	require.True(t, validate.IsValidationError(wrap))
}

func TestLastSegment(t *testing.T) {
	err := validate.Field(
		"address.name",
		"",
		failValidatorWithCode[string]("fail"),
	)

	errs := validate.Collect(err)
	require.Equal(t, 1, len(errs))
	require.Equal(t, "address.name", errs[0].Path)
	require.Equal(t, "name", validate.LastPathSegment(errs[0].Path))
}

func TestMapError(t *testing.T) {
	type Person struct {
		Name    string
		Address string
	}

	t.Run("map values", func(t *testing.T) {
		err := validate.Map(
			"persons",
			map[string]Person{
				"some-id": {
					Name:    "John",
					Address: "Street 1",
				},
				"multiple-some-id": {
					Name:    "John",
					Address: "Street 1",
				},
			},
		).Values("person", func(value Person) error {
			return validate.Join(
				validate.Field("name", value.Name, failValidatorWithCode[string]("fail")),
				validate.Field("address", value.Address, failValidatorWithCode[string]("fail")),
			)
		})
		require.NotNil(t, err)
		errs := validate.Collect(err)
		require.Equal(t, 4, len(errs))

		for _, err := range errs {
			if !slices.Contains([]string{"persons.some-id.person.name", "persons.some-id.person.address", "persons.multiple-some-id.person.name", "persons.multiple-some-id.person.address"}, err.ExactPath) {
				t.Errorf("unexpected exact path: %s", err.ExactPath)
			}
		}
	})

	t.Run("map keys", func(t *testing.T) {
		err := validate.Map(
			"persons",
			map[string]Person{
				"some-id": Person{
					Name:    "John",
					Address: "Street 1",
				},
			},
		).Keys("person", func(key string) error {
			return validate.Field("name", key, failValidatorWithCode[string]("fail"))
		})
		require.NotNil(t, err)
		errs := validate.Collect(err)
		require.Equal(t, 1, len(errs))
		require.Equal(t, "persons.person.name", errs[0].Path)
		require.Equal(t, "persons.some-id.person.name", errs[0].ExactPath)
	})

	t.Run("map key", func(t *testing.T) {
		err := validate.Map(
			"persons",
			map[string]Person{
				"some-id": Person{
					Name:    "John",
					Address: "Street 1",
				},
			},
		).Key("person", "some-id", func(value Person) error {
			return validate.Join(
				validate.Field("name", value.Name, failValidatorWithCode[string]("fail")),
				validate.Field("address", value.Address, failValidatorWithCode[string]("fail")),
			)
		})
		require.NotNil(t, err)
		errs := validate.Collect(err)
		require.Equal(t, 2, len(errs))
		require.Equal(t, "persons.person.name", errs[0].Path)
		require.Equal(t, "persons.some-id.person.name", errs[0].ExactPath)
		require.Equal(t, "persons.person.address", errs[1].Path)
		require.Equal(t, "persons.some-id.person.address", errs[1].ExactPath)
	})

	t.Run("map slice items", func(t *testing.T) {
		list := []Person{
			{
				Name:    "John",
				Address: "Street 1",
			},
		}

		err := validate.Slice("persons", list).Items("person", func(value Person) error {
			return validate.Join(
				validate.Field("name", value.Name, failValidatorWithCode[string]("fail")),
				validate.Field("address", value.Address, failValidatorWithCode[string]("fail")),
			)
		})
		require.NotNil(t, err)
		errs := validate.Collect(err)
		require.Equal(t, 2, len(errs))
		require.Equal(t, "persons.*.person.name", errs[0].Path)
		require.Equal(t, "persons.0.person.name", errs[0].ExactPath)
		require.Equal(t, "persons.*.person.address", errs[1].Path)
		require.Equal(t, "persons.0.person.address", errs[1].ExactPath)
	})
}
