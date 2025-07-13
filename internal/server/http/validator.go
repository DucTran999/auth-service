package http

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// registerStrongPasswordValidators registers a custom 'password' validation rule
// that enforces strong password requirements.
func RegisterStrongPasswordValidators(v *validator.Validate) error {
	return v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()

		if len(password) < 8 {
			return false
		}

		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString
		hasLower := regexp.MustCompile(`[a-z]`).MatchString
		hasDigit := regexp.MustCompile(`[0-9]`).MatchString
		hasSymbol := regexp.MustCompile(`[\W_]`).MatchString

		return hasUpper(password) && hasLower(password) && hasDigit(password) && hasSymbol(password)
	})
}
