package utils

import "unicode"

func ValidatePassword(str string) bool {
	var hasUppercase, hasLowercase, hasDigits, hasSpecialChar bool

	if len(str) < 8 || len(str) > 128 {
		return false
	}

	for _, c := range str {
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

	if hasUppercase || hasLowercase || hasDigits || hasSpecialChar {
		return true
	}

	return false
}
