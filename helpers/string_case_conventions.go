package helpers

import (
	"strings"
)

func ToLowerCamelCase(s string) string {
	if len(s) <= 1 {
		return strings.ToLower(s)
	} else {
		return strings.ToLower(s[0:1]) + s[1:]
	}
}

func prefixCapitalLettersWithString(origStr, prefixWith string) string {
	finalStr := ""
	for i, _ := range origStr {
		currentLetterAsStr := origStr[i : i+1]

		if i > 0 && strings.ToLower(currentLetterAsStr) != currentLetterAsStr {
			finalStr += prefixWith
		}
		finalStr += currentLetterAsStr
	}
	return finalStr
}

func ToKebabCase(s string) string {
	return strings.ToLower(prefixCapitalLettersWithString(s, "-"))
}

func ToSnakeCase(s string) string {
	return strings.ToLower(prefixCapitalLettersWithString(s, "_"))
}
