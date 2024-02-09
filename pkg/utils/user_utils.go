package utils

import (
	"regexp"
	"strings"

	"org.com/org/pkg/database/mongodb/repository"
)

func IsValidEmail(email string) bool {
	// Regular expression for a simple email validation
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}

func ValidateEmail(email string) string {
	// Check email format
	if !IsValidEmail(email) {
		return "Invalid email format"
	}

	// Check if email already exists
	isEmailExist, err := repository.IsEmailExists(email)
	if err != nil {
		return "Error checking email existence"
	}
	if isEmailExist {
		return "This Email already exists, please use another one"
	}

	// All validations passed
	return ""
}

func ExtractErrorMessage(err error) string {
    // Extracting the specific error message
    parts := strings.Split(err.Error(), "Error:Field validation for ")
    if len(parts) > 1 {
        return parts[1]
    }
    return err.Error()
}