package validator

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	V *validator.Validate
	T *ut.Translator
}

func NewValidation() *Validator {
	v := &Validator{
		V: validator.New(),
	}
	trans := initializeTranslation(v.V)
	v.T = trans
	registerFunc(v.V)
	return v
}

func initializeTranslation(validate *validator.Validate) *ut.Translator {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)
	return &trans
}

func registerFunc(validate *validator.Validate) {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func (v *Validator) Validate(form interface{}) []error {
	var errResp []error
	if err := v.V.Struct(form); err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			err := errors.New(e.Translate(*v.T))
			// err := errors.New(e.Translate(*v.T), &errors.BadRequest)
			// key := strings.SplitAfterN(e.Namespace(), ".", 2)
			// err = errors.SetContext(err, key[1], e.Translate(*v.T))
			errResp = append(errResp, err)
		}
	}
	return errResp
}
