package fiberx

import (
	"regexp"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ent "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2/utils"
)

// V is an instance of *validator.Validate
// used for custom registrations.
var V *validator.Validate

var (
	uni   *ut.UniversalTranslator
	trans ut.Translator

	mobileRegexp *regexp.Regexp
)

func init() {
	V = validator.New()

	uni = ut.New(en.New())
	trans, _ = uni.GetTranslator("en")

	if err := ent.RegisterDefaultTranslations(V, trans); err != nil {
		panic(err)
	}

	registerValidations()
}

func registerValidations() {
	mobileRegexp = regexp.MustCompile(`^1[3456789]\d{9}$`)

	_ = V.RegisterValidation("mobile", MobileRule)
}

// MobileRule validates mobile phone number
func MobileRule(fl validator.FieldLevel) bool {
	return mobileRegexp.Match(utils.GetBytes(fl.Field().String()))
}
