package validator

import (
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	alphaSpaceRegex = regexp.MustCompile(`^[a-zA-Z\s]+$`)
	usernameRegex   = regexp.MustCompile(`^[a-zA-Z0-9@._]{5,30}$`)
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
	validate.RegisterValidation("alpha_space", func(fl validator.FieldLevel) bool {
		return alphaSpaceRegex.MatchString(fl.Field().String())
	})
	validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		return usernameRegex.MatchString(fl.Field().String())
	})
	validate.RegisterValidation("password", passwordValidator)
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}
	allowedChars := regexp.MustCompile(`^[A-Za-z\d@$!%*?&]+$`)
	if !allowedChars.MatchString(password) {
		return false
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return false
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false
	}
	if !regexp.MustCompile(`\d`).MatchString(password) {
		return false
	}
	if !regexp.MustCompile(`[@$!%*?&]`).MatchString(password) {
		return false
	}
	return true
}

func (v *Validator) Validate(form interface{}) []error {
	var errResp []error
	if err := v.V.Struct(form); err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			if e.Tag() == "password" {
				err := errors.New("weak password")
				errResp = append(errResp, err)
				continue
			}
			if e.Tag() == "username" {
				err := errors.New("invalid username")
				errResp = append(errResp, err)
				continue
			}
			err := errors.New(e.Translate(*v.T))
			errResp = append(errResp, err)
		}
	}
	return errResp
}
