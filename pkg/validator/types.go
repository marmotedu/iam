package validator

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Struct(data interface{}) error
	StructCtx(ctx context.Context, data interface{}) error
	Var(f interface{}, rule string) error
	VarCtx(ctx context.Context, f interface{}, rule string) error
	WithTranslator(language string) Validator
	GetValidate() *validator.Validate
	RegisterValidation(tag, errMSg string, vf Validation) error
}

// Validation custom validation func.
type Validation func(fl validator.FieldLevel) bool
