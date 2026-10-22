package utils

import "unicode"

func ValidatePasswordPolicy(password string) bool {
	var hasUppercase, hasLowercase, hasDigits, hasSpecialChar bool

	if len(password) < 8 || len(password) > 128 {
		return false
	}

	for _, c := range password {
		switch {
		case unicode.IsLower(c):
			hasLowercase = true
		case unicode.IsUpper(c):
			hasUppercase = true
		case unicode.IsDigit(c):
			hasDigits = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecialChar = true
		}
	}

	if hasUppercase && hasLowercase && hasDigits && hasSpecialChar {
		return true
	}

	return false
}
