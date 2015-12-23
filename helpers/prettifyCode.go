package helpers

import (
	"strings"
)

type PredicateFunc func(trimmedLine string) bool

type PrettifyRules struct {
	MustPrefixWithEmptyLine  PredicateFunc
	StartIndentNextLine      PredicateFunc
	StopIndentingCurrentLine PredicateFunc
}

func PrettifyCode(code []byte, rules *PrettifyRules) []byte {
	originalLines := strings.Split(string(code), "\n")
	prettyLines := []string{}

	currentIndent := ""
	indentStr := "\t"
	for _, line := range originalLines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}

		linePrefix := ""
		if rules != nil && rules.MustPrefixWithEmptyLine != nil && rules.MustPrefixWithEmptyLine(trimmedLine) {
			linePrefix = "\n"
		}

		if rules != nil && rules.StopIndentingCurrentLine != nil && rules.StopIndentingCurrentLine(trimmedLine) {
			currentIndent = currentIndent[:len(currentIndent)-len(indentStr)]
		}

		prettyLines = append(prettyLines, linePrefix+currentIndent+trimmedLine)

		if rules != nil && rules.StartIndentNextLine != nil && rules.StartIndentNextLine(trimmedLine) {
			currentIndent += indentStr
		}
	}

	return []byte(strings.Join(prettyLines, "\n"))
}
