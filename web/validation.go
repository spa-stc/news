package web

import (
	"errors"
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	validator  *validator.Validate
	translator ut.Translator
}

func NewValidator() *Validator {
	v := validator.New()
	en := en.New()
	uni := ut.New(en, en)

	trans, _ := uni.GetTranslator("en")

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		panic(err)
	}

	return &Validator{
		validator:  v,
		translator: trans,
	}
}

func (v *Validator) Struct(data any) (map[string]string, error) {
	err := v.validator.Struct(data)
	if err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			return errs.Translate(v.translator), nil
		}

		return nil, fmt.Errorf("error running validator: %w", err)
	}

	return nil, nil //nolint:nilnil // This is the api I want to present.
}
