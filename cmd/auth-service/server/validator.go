package server

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func registerStrongPasswordValidators(v *validator.Validate) {
	// Register custom 'password' validation
	v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()

		if len(password) < 8 {
			return false
		}

		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSymbol := regexp.MustCompile(`[\W_]`).MatchString(password)

		return hasUpper && hasLower && hasNumber && hasSymbol
	})
}
