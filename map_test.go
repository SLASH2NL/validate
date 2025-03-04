package validate_test

import (
	"testing"

	"github.com/SLASH2NL/validate"
	"github.com/stretchr/testify/require"
)

func TestValidateMapKey(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		data := map[string]string{
			"somekey": "value",
		}

		err := validate.Map("items", data).Key("name", "somekey", validate.MinString(10))
		require.Error(t, err)
		errs := validate.Collect(err)
		require.Equal(t, 1, len(errs))
		require.Equal(t, "items.name", errs[0].Path)
		require.Equal(t, "items.somekey.name", errs[0].ExactPath)
		require.Equal(t, validate.CodeStringMin, errs[0].Violations[0].Code)
		require.Equal(t, "somekey", errs[0].Args["key"])
	})

	t.Run("missing key", func(t *testing.T) {
		data := map[string]string{
			"somekey": "value",
		}

		err := validate.Map("items", data).Key("email", "missingkey", validate.MinString(10))
		require.Error(t, err)
		errs := validate.Collect(err)
		require.Equal(t, 1, len(errs))
		require.Equal(t, "items.email", errs[0].Path)
		require.Equal(t, "items.missingkey.email", errs[0].ExactPath)
		require.Equal(t, validate.CodeUnknownField, errs[0].Violations[0].Code)
		require.Equal(t, "missingkey", errs[0].Args["key"])
	})

	t.Run("resolve", func(t *testing.T) {
		type d struct {
			Name string
		}

		data := map[string]d{
			"somekey": {
				Name: "invalid",
			},
		}

		err := validate.Map("items", data).Key("name", "somekey", validate.Resolve(func(t d) string { return t.Name }, validate.MinString(10))...)
		require.Error(t, err)
		errs := validate.Collect(err)
		require.Equal(t, 1, len(errs))
		require.Equal(t, "items.name", errs[0].Path)
		require.Equal(t, "items.somekey.name", errs[0].ExactPath)
		require.Equal(t, validate.CodeStringMin, errs[0].Violations[0].Code)
	})

	t.Run("no error", func(t *testing.T) {
		data := map[string]string{
			"somekey": "value",
		}

		err := validate.Map("items", data).Key("email", "somekey", validate.MinString(5))
		require.NoError(t, err)
	})
}

func TestValidateMapKeys(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		data := map[string]int{
			"somekey":            5,
			"otherthatiscorrect": 10,
		}

		err := validate.Map("items", data).Keys("name", validate.MinString(10))
		require.Error(t, err)
		errs := validate.Collect(err)
		require.Equal(t, 1, len(errs))
		require.Equal(t, "items.name", errs[0].Path)
		require.Equal(t, "items.somekey.name", errs[0].ExactPath)
		require.Equal(t, validate.CodeStringMin, errs[0].Violations[0].Code)
		require.Equal(t, "somekey", errs[0].Args["key"])
	})

	t.Run("no error", func(t *testing.T) {
		data := map[string]int{
			"somekey": 5,
			"other":   10,
		}

		err := validate.Map("items", data).Keys("email", validate.MinString(5))
		require.NoError(t, err)
	})
}

func TestValidateMapValues(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		data := map[string]string{
			"somekey": "test@example.org",
			"other":   "not-email",
		}

		err := validate.Map("items", data).Values("name", validate.Email)
		require.Error(t, err)
		errs := validate.Collect(err)
		require.Equal(t, 1, len(errs))
		require.Equal(t, "items.name", errs[0].Path)
		require.Equal(t, "items.other.name", errs[0].ExactPath)
		require.Equal(t, validate.CodeEmail, errs[0].Violations[0].Code)
		require.Equal(t, "other", errs[0].Args["key"])
	})

	t.Run("resolve", func(t *testing.T) {
		type d struct {
			Email string
		}

		data := map[string]d{
			"somekey": {
				Email: "test@example.org",
			},
			"other": {
				Email: "not-email",
			},
		}

		err := validate.Map("items", data).Values("email", validate.Resolve(func(t d) string { return t.Email }, validate.Email)...)
		require.Error(t, err)
		errs := validate.Collect(err)
		require.Equal(t, 1, len(errs))
		require.Equal(t, "items.email", errs[0].Path)
		require.Equal(t, "items.other.email", errs[0].ExactPath)
		require.Equal(t, validate.CodeEmail, errs[0].Violations[0].Code)
		require.Equal(t, "other", errs[0].Args["key"])
	})

	t.Run("no error", func(t *testing.T) {
		data := map[string]string{
			"somekey": "test@example.org",
			"other":   "test2@example.org",
		}

		err := validate.Map("items", data).Values("email", validate.Email)
		require.NoError(t, err)
	})
}
