package validator

import (
	"context"
	"reflect"
	"sync"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	once sync.Once
	mu   sync.Mutex
	v    Validator
)

func init() {
	once.Do(func() {
		opts := NewOptions()
		v = New(opts.Language, opts.Tag)
	})
}

func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	v = New(opts.Language, opts.Tag)
}

func New(language, tag string) Validator {
	v := &validate{
		Language: language,
		Tag:      tag,
		Validate: validator.New(),
	}
	v.WithTranslator(v.Language)
	v.registerTagName()
	return v
}

func Struct(data interface{}) error {
	return v.Struct(data)
}

func StructCtx(ctx context.Context, data interface{}) error {
	return v.StructCtx(ctx, data)
}

func Var(f interface{}, rule string) error {
	return v.Var(f, rule)
}

func VarCtx(ctx context.Context, f interface{}, rule string) error {
	return v.VarCtx(ctx, f, rule)
}

// WithTranslator set default translation.
func WithTranslator(language string) Validator {
	return v.WithTranslator(language)
}

func GetValidate() *validator.Validate {
	return v.GetValidate()
}

// RegisterValidation register custom validation func.
func RegisterValidation(tag, errMsg string, vf Validation) error {
	return v.RegisterValidation(tag, errMsg, vf)
}

type validate struct {
	Language string
	Tag      string
	trans    ut.Translator
	*validator.Validate
}

// WithTranslator set translation.
func (v *validate) WithTranslator(language string) Validator {
	v.Language = language
	var (
		zhT = zh.New() // Chinese trans.
		enT = en.New() // English trans.
		// If you need more language.
	)
	uni := ut.New(zhT, enT)
	v.trans, _ = uni.GetTranslator(v.Language)

	switch v.Language {
	case "zh":
		_ = zh_translations.RegisterDefaultTranslations(v.Validate, v.trans)
	case "en":
		_ = en_translations.RegisterDefaultTranslations(v.Validate, v.trans)
	}

	return v
}

func (v *validate) Struct(data interface{}) error {
	if err := v.Validate.Struct(data); err != nil {
		return &ValidationErrors{
			trans: v.trans,
			errs:  err.(validator.ValidationErrors),
		}
	}

	return nil
}

func (v *validate) StructCtx(ctx context.Context, data interface{}) error {
	if err := v.Validate.StructCtx(ctx, data); err != nil {
		return &ValidationErrors{
			trans: v.trans,
			errs:  err.(validator.ValidationErrors),
		}
	}

	return nil
}

func (v *validate) Var(f interface{}, rule string) error {
	if err := v.Validate.Var(f, rule); err != nil {
		return &ValidationErrors{
			trans: v.trans,
			errs:  err.(validator.ValidationErrors),
		}
	}

	return nil
}

func (v *validate) VarCtx(ctx context.Context, f interface{}, rule string) error {
	if err := v.Validate.VarCtx(ctx, f, rule); err != nil {
		return &ValidationErrors{
			trans: v.trans,
			errs:  err.(validator.ValidationErrors),
		}
	}

	return nil
}

// RegisterValidation register custom validation func.
func (v *validate) RegisterValidation(tag, errMsg string, vf Validation) error {
	if err := v.Validate.RegisterValidation(tag, validator.Func(vf)); err != nil {
		return err
	}

	if err := v.registerValidationTranslator(tag, errMsg); err != nil {
		return err
	}

	return nil
}

func (v *validate) GetValidate() *validator.Validate {
	return v.Validate
}

// Register struct field tag name.
func (v *validate) registerTagName() {
	v.Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get(v.Tag)
	})
}

// register custom validation translator.
func (v *validate) registerValidationTranslator(tag string, msg string) error {
	f := func(ut ut.Translator) error {
		if err := v.trans.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}

	return v.Validate.RegisterTranslation(tag, v.trans, f, translate)
}

func translate(trans ut.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		panic(fe.(error).Error())
	}
	return msg
}
