package validate_test

import (
	"fmt"

	"github.com/SLASH2NL/validate"
)

func ExampleValidate() {
	o := Order{
		Customer: Customer{
			FirstName: "",
			LastName:  "Doe's last name is too long",
			Email:     "john@doe.com",
		},
		ShippingAddress: Address{
			Street: "Main Street",
			Number: 123,
		},
		BillingAddress: Address{
			Street: "Second Street is way too long",
			Number: 321,
		},
		Products: []Product{
			{
				Name:  "Product 1",
				Price: -1,
				Qty:   1,
			},
			{
				Name:  "Product 2",
				Price: 200,
				Qty:   -2,
			},
			{
				Name:  "Product 3 has a name that is too long",
				Price: 300,
				Qty:   3,
			},
		},
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}

	err := validate.Join(
		validate.Group(
			"customer",
			validate.Join(
				validate.Field(
					"first_name",
					o.Customer.FirstName,
					validate.Required,
				),
				validate.Field(
					"last_name",
					o.Customer.LastName,
					validate.Required,
					validate.StrMax(15),
				),
			),
		),
		validate.Group(
			"shipping_address",
			o.ShippingAddress.validate(),
		),
		validate.Group(
			"billing_address",
			o.BillingAddress.validate(),
		),
		validate.Group(
			"products",
			validate.Slice(
				o.Products,
				func(p Product) error {
					return validate.Join(
						validate.Field(
							"name",
							p.Name,
							validate.Required,
							validate.StrMax(30),
						),
						validate.Field(
							"price",
							p.Price,
							validate.NumberMin(0),
						),
						validate.Field(
							"qty",
							p.Qty,
							validate.NumberMin(1),
						),
					)
				},
			),
		),
		validate.Group(
			"options",
			validate.Map(
				o.Options,
				validate.Key("option1", validate.OneOf("accepted", "values")),
				validate.Key("option3", validate.Required[string]),
			),
		),
	)

	errs := validate.Collect(err)
	for _, e := range errs {
		fmt.Printf("Path: %s, Field: %s, Code: %s, Args: %v\n", e.Path, e.Field, e.Code, e.Args)
	}

	// Output:
	// Path: customer, Field: first_name, Code: required, Args: map[]
	// Path: customer, Field: last_name, Code: string.max, Args: map[max:15]
	// Path: billing_address, Field: street, Code: string.max, Args: map[max:15]
	// Path: products.[0], Field: price, Code: number.min, Args: map[min:0]
	// Path: products.[1], Field: qty, Code: number.min, Args: map[min:1]
	// Path: products.[2], Field: name, Code: string.max, Args: map[max:30]
	// Path: options, Field: option1, Code: oneof, Args: map[accepted:[accepted values]]
	// Path: options, Field: option3, Code: unknown, Args: map[]
}

type Order struct {
	Customer        Customer
	ShippingAddress Address
	BillingAddress  Address
	Products        []Product
	Options         map[string]string
}

type Customer struct {
	FirstName string
	LastName  string
	Email     string
}

type Address struct {
	Street string
	Number int
}

func (a Address) validate() error {
	return validate.Join(
		validate.Field(
			"street",
			a.Street,
			validate.Required,
			validate.StrMax(15),
		),
		validate.Field(
			"number",
			a.Number,
			validate.Required,
		),
	)
}

type Product struct {
	Name  string
	Price int
	Qty   int
}

// func (x SomeStruct) validate() error {
// 	return validate.Join(
// 		validate.Field(
// 			"first_name",
// 			x.FirstName,
// 			validate.Required,
// 		),
// 		validate.Field(
// 			"last_name",
// 			x.LastName,
// 			validate.FailFirst(
// 				validate.Required,
// 				validate.StrMax(255),
// 			),
// 		),
// 		// You can group validators to validate nested data types.
// 		validate.Group("address", validate.Validate(
// 			"street",
// 			x.Address.Street,
// 			validate.Required,
// 		)),
// 		// Validate a map.
// 		validate.Group("args", validate.Map(
// 			x.Args,
// 			validate.Key("location", validate.StrMax(255)),
// 		)),
// 		// Validate a slice.
// 		validate.Group("notes", validate.Slice(
// 			x.Notes,
// 			validate.GroupValidators( // Use GroupValidators to create a proper field name. If not used the field name will be just the index of the offending item like: [0] or [1].
// 				"text",
// 				x.Notes.Text,
// 				validate.Required,
// 			),
// 		)),
// 	)
// }
