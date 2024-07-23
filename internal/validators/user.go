package validators

import (
	"fmt"
	"net/mail"
	"unicode"
)

// ValidateUsername validate username
func ValidateUsername(username string) error {
	// Maximum 20 and minimum 1
	if len(username) > 20 || len(username) == 0 {
		return fmt.Errorf("length username must be less 20 and not empty")
	}

	// only lowercase
	for _, r := range username {
		if unicode.IsUpper(r) && unicode.IsLetter(r) {
			return fmt.Errorf("username only lowercase")
		}
	}

	// letters, numbers, underscore (_) and period (.)
	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '.' && r != '_' {
			return fmt.Errorf("username must be letters, numbers, underscore (_) and period (.)")
		}
	}

	// no period/point in the end allowed
	lastCharacter := username[len(username)-1:]
	if lastCharacter == "." {
		return fmt.Errorf("lastCharacter must not be period")
	}

	return nil
}

// ValidateEmail validate email
func ValidateEmail(email string) error {
	// TO DO e.g. email "Barry Gibbs <bg@example.com>" return err nil
	_, err := mail.ParseAddress(email)
	return err
}

// ValidatePassword validate password
func ValidatePassword(password string) error {
	// minimum 8
	if len(password) < 8 {
		return fmt.Errorf("length password must be larger 8")
	}

	// upper case
	upper := false
	for _, r := range password {
		if unicode.IsUpper(r) && unicode.IsLetter(r) {
			upper = true
			break
		}
	}

	// lower case
	lower := false
	for _, r := range password {
		if unicode.IsLower(r) && unicode.IsLetter(r) {
			lower = true
			break
		}
	}

	// number
	number := false
	for _, r := range password {
		if unicode.IsNumber(r) {
			number = true
			break
		}
	}
	if !upper || !number || !lower {
		return fmt.Errorf("password must contain upper, lower, number")
	}

	return nil
}

// ValidateFIO validate user first name and last name
func ValidateFIO(value string) error {
	// Maximum 20 and minimum 1
	if len(value) > 30 || len(value) == 0 {
		return fmt.Errorf("length fio must be less 30 and not empty")
	}

	// letters, numbers, underscore (_) and period (.)
	for _, r := range value {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '.' && r != '_' {
			return fmt.Errorf("fio must be letters, numbers, underscore (_) and period (.)")
		}
	}

	// no period/point in the end allowed
	lastCharacter := value[len(value)-1:]
	if lastCharacter == "." {
		return fmt.Errorf("lastCharacter must not be period")
	}

	return nil
}
