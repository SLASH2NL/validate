package validate_test

import (
	"fmt"

	"github.com/SLASH2NL/validate"
)

func ExampleField() {
	err := validate.Field("email", "test", validate.Email)
	printError(err)

	// Output:
	// validation error for exact path: email, path: email, args: map[], violations: [{email map[]}]
}

func ExampleJoin() {
	err := validate.Join(
		validate.Field("email", "test", validate.Email),
		validate.Field("iban", "invalid", validate.IBAN),
	)
	printError(err)

	// Output:
	// validation error for exact path: email, path: email, args: map[], violations: [{email map[]}]
	// validation error for exact path: iban, path: iban, args: map[], violations: [{iban map[]}]
}

func ExampleCollect() {
	err := validate.Join(
		validate.Field("email", "test", validate.Email),
		validate.Field("iban", "invalid", validate.IBAN),
	)

	errs := validate.Collect(err)
	fmt.Printf("Size is %d", len(errs))

	// Output:
	// Size is 2
}

func ExampleIf() {
	err := validate.Join(
		validate.Field("email", "test", validate.If(true, validate.Email)...),
		validate.Field("iban", "invalid", validate.If(false, validate.IBAN)...),
	)
	printError(err)

	// Output:
	// validation error for exact path: email, path: email, args: map[], violations: [{email map[]}]
}

func ExampleFailFirst() {
	err := validate.Field("email", "test", validate.FailFirst(validate.Email, validate.MinString(100)))
	printError(err)

	// Output:
	// validation error for exact path: email, path: email, args: map[], violations: [{email map[]}]
}

func ExampleOverridePath() {
	err := validate.OverridePath("name", validate.Field("email", "test", validate.Email, validate.MinString(100)))
	printError(err)

	// Output:
	// validation error for exact path: email, path: name, args: map[], violations: [{email map[]} {min.string map[min:100]}]
}

func ExampleOverrideExactPath() {
	err := validate.OverrideExactPath("name", validate.Field("email", "test", validate.Email, validate.MinString(100)))
	printError(err)

	// Output:
	// validation error for exact path: name, path: email, args: map[], violations: [{email map[]} {min.string map[min:100]}]
}

func ExampleSlice() {
	data := []testSlice{
		{Name: "John Deer", Amount: 9},
		{Name: "Deer John", Amount: 1},
	}

	err := validate.Slice("data", data).Items("total", validate.Resolve(func(value testSlice) int { return value.Amount }, validate.MaxNumber(5))...)
	printError(err)

	// Output:
	// validation error for exact path: data.0.total, path: data.*.total, args: map[index:0], violations: [{max.number map[max:5]}]
}

func ExampleMap() {
	data := map[string]string{
		"first":  "john deer",
		"second": "Deer John",
	}

	err := validate.Map("data", data).Values("name", validate.Lowercase)
	printError(err)

	// Output:
	// validation error for exact path: data.name, path: data.name, args: map[key:second], violations: [{lowercase map[]}]
}

func printError(err error) {
	if err == nil {
		return
	}

	errs := validate.Collect(err)
	for _, e := range errs {
		fmt.Println(e.Error())
	}
}
