package handler

import "github.com/go-playground/validator/v10"

func validationErrorMessage(fe validator.FieldError) string {
	switch fe.Field() {
	case "Email":
		switch fe.Tag() {
		case "required":
			return "Email is required"
		case "email":
			return "Email must be valid"
		}
	case "Password":
		switch fe.Tag() {
		case "required":
			return "Password is required"
		case "password":
			return "Password must include at least 1 uppercase, 1 lowercase, 1 number, and 1 special character"
		}
	}

	return "Invalid input"
}
