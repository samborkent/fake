package cases

import "unicode"

// Camel converts Pascal or upper camel case strings to lower camel case.
func Camel(str string) string {
	if len(str) == 0 {
		return str
	}

	runes := []rune(str)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// Snake converts camel or Pascal case strings to snake case.
func Snake(str string) string {
	if len(str) == 0 {
		return str
	}

	var result []rune

	runes := []rune(str)

	result = append(result, unicode.ToLower(runes[0]))

	for i := 1; i < len(runes); i++ {
		if unicode.IsUpper(runes[i]) {
			if !unicode.IsUpper(runes[i-1]) {
				result = append(result, '_')
			}

			result = append(result, unicode.ToLower(runes[i]))
		} else {
			result = append(result, runes[i])
		}
	}

	return string(result)
}
