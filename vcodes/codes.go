package vcodes

type Code string

const (
	Unknown            Code = "unknown_error"
	NotFound           Code = "not.found"
	Required           Code = "required"
	Equal              Code = "equal"
	OneOf              Code = "oneof"
	NumberMin          Code = "number.min"
	NumberMax          Code = "number.max"
	StringEmail        Code = "string.email"
	StringRegex        Code = "string.regex"
	StringRegexInvalid Code = "string.regex.invalid"
	StringMin          Code = "string.min"
	StringMax          Code = "string.max"
	StringLowercase    Code = "string.lowercase"
	StringUppercase    Code = "string.uppercase"
	IBAN               Code = "iban"
)
