package validator

import (
	"bytes"
	"fmt"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ValidationErrors struct {
	trans ut.Translator
	errs  validator.ValidationErrors
}

func (e *ValidationErrors) TranslateErrs() (errs []error) {
	for _, val := range e.errs {
		errs = append(errs, &FieldError{
			Field:   val.Field(),
			Message: val.Translate(e.trans),
		})
	}

	return
}

func (e *ValidationErrors) TranslateErrsMap() map[string]string {
	return removeTopStruct(e.errs.Translate(e.trans))
}

func (e *ValidationErrors) GetValidatorValidationErrors() validator.ValidationErrors {
	return e.errs
}

func (e *ValidationErrors) Error() string {
	buff := bytes.NewBufferString("")
	for _, val := range e.errs {
		buff.WriteString(val.Translate(e.trans))
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

type FieldError struct {
	Field   string
	Message string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}

	return res
}
