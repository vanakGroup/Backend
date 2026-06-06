package utils

import (
	"regexp"
	"strings"
)

func IsNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func CheckLicensePlate(LicensePlate string) (*string, bool) {
	invalidLetters := map[string]bool{
		"ث": true, "ح": true, "خ": true, "ر": true, "ز": true,
		"ض": true, "ظ": true, "غ": true, "ف": true, "ک": true,
		"ی": true,
	}

	normalizer := strings.NewReplacer(
		"۰", "0", "۱", "1", "۲", "2", "۳", "3", "۴", "4",
		"۵", "5", "۶", "6", "۷", "7", "۸", "8", "۹", "9",
		"٠", "0", "١", "1", "٢", "2", "٣", "3", "٤", "4",
		"٥", "5", "٦", "6", "٧", "7", "٨", "8", "٩", "9",
		"ي", "ی", "ك", "ک",
	)

	cleaned := normalizer.Replace(LicensePlate)

	plateRegex := regexp.MustCompile(`^([1-9]\d)([ب-ی]|الف)(\d{3})([1-9]\d)$`)
	matches := plateRegex.FindStringSubmatch(cleaned)

	if len(matches) != 5 {
		return nil, false
	}

	letter := matches[2]

	if invalidLetters[letter] {
		return nil, false
	}

	return &cleaned, true
}
