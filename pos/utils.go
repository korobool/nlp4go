package pos

func IsDigit(s string) bool {
	for _, char := range s {
		if !(char >= '0' && char <= '9') {
			return false
		}
	}
	return true
}

func StrSuffix(str string, sz int) string {
	rstr := []rune(str)
	if len(rstr) < sz {
		return str
	}
	return string(rstr[len(rstr)-sz:])
}

func StrPrefix(str string, sz int) string {
	rstr := []rune(str)
	if len(rstr) < sz {
		return str
	}
	return string(rstr[:len(rstr)-sz])
}
