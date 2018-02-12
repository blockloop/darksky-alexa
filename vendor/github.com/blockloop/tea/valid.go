package tea

import (
	"context"

	"gopkg.in/go-playground/validator.v9"
)

// Validate is the Validator used to validate structs within Body
var Validate Validator = validator.New()

// NopValidator is a validator that always returns nil and does not
// do any validation
var NopValidator Validator = &nopValidator{}

// Validator is a struct validator used within Body
type Validator interface {
	StructCtx(context.Context, interface{}) error
	VarCtx(context.Context, interface{}, string) error
}

type nopValidator struct{}

func (*nopValidator) StructCtx(context.Context, interface{}) error {
	return nil
}
func (*nopValidator) VarCtx(context.Context, interface{}, string) error {
	return nil
}
