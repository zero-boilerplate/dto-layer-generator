package helpers

import (
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"github.com/francoishill/golang-web-dry/osutils"
	"io/ioutil"
	"strings"
)

func replaceTabsWithSpaces(s string) string {
	return strings.Replace(s, "\t", "    ", -1)
}

func InjectContentIntoFilePlaceholder(logger Logger, filePath, placeholder, content string) {
	finalContent := content
	globalSnippetIndentSpaces := ""
	if osutils.FileExists(filePath) {
		logger.Debug("Found existing file '%s'. Replacing the placeholders.", filePath)

		fileContent, err := ioutil.ReadFile(filePath)
		CheckError(err)

		foundPlaceholder := false
		isBeginPlaceholder := false

		beginText := "{{BEGIN " + placeholder + "}}"
		endText := "{{END " + placeholder + "}}"

		origLines := strings.Split(string(fileContent), "\n")
		finalLines := []string{}
		for _, originalLine := range origLines {
			trimmedLine := strings.TrimSpace(originalLine)
			lineWithTabsReplacedBySpaces := replaceTabsWithSpaces(originalLine)

			if strings.Contains(trimmedLine, endText) {
				isBeginPlaceholder = false
				finalLines = append(finalLines, globalSnippetIndentSpaces+trimmedLine)
				continue
			}

			if strings.Contains(trimmedLine, beginText) {
				isBeginPlaceholder = true
				foundPlaceholder = true

				//Indent all the code with the same indentation as the 'begin' placeholder line
				prefixSpacesCount := len(lineWithTabsReplacedBySpaces) - len(strings.TrimLeft(lineWithTabsReplacedBySpaces, " "))
				globalSnippetIndentSpaces = strings.Repeat(" ", prefixSpacesCount)

				finalLines = append(finalLines, globalSnippetIndentSpaces+trimmedLine)

				contentLines := strings.Split(content, "\n")
				for index, _ := range contentLines {
					contentLines[index] = globalSnippetIndentSpaces + replaceTabsWithSpaces(contentLines[index])
				}
				finalLines = append(finalLines, strings.Join(contentLines, "\n"))

				logger.Debug("Appended content after line containing '%s'", beginText)

				continue
			}

			if isBeginPlaceholder {
				//Skip when between the placeholder tags
				continue
			}

			//These are the code lines "outside" of the placeholder block, so keep the lines untouched
			finalLines = append(finalLines, originalLine)
		}

		finalContent = strings.Join(finalLines, "\n")

		if !foundPlaceholder {
			logger.Error("Placeholder '%s' not found in '%s'", placeholder, filePath)
		}
	}
	err := ioutil.WriteFile(filePath, []byte(finalContent), 0600)
	CheckError(err)
}
