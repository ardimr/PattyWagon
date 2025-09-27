package utils

import (
	"net/mail"
	"regexp"
	"strings"
)

func IsEmail(email string) bool {
	emailAddress, err := mail.ParseAddress(email)
	return err == nil && emailAddress.Address == email
}

func IsValidPhoneNumber(phone string) bool {
	if phone == "" {
		return false
	}

	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	phone = strings.TrimSpace(phone)

	phoneRegex := regexp.MustCompile(`^\+[1-9]\d{6,14}$`)
	return phoneRegex.MatchString(phone)
}
