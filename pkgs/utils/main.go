package utils

func InArray(s string, arr []string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}

	return false
}

func DerefString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
