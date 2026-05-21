package utils

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
