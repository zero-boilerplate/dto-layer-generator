package helpers

import (
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"github.com/francoishill/golang-web-dry/osutils"
	"io/ioutil"
	"strings"
)

func InjectContentIntoFilePlaceholder(logger Logger, filePath, placeholder, content string) {
	finalContent := content
	globalSnippetIndentSpaces := ""
	if osutils.FileExists(filePath) {
		logger.Debug("Found existing file '%s'. Replacing the placeholders.", filePath)

		fileContent, err := ioutil.ReadFile(filePath)
		CheckError(err)

		isBeginPlaceholder := false

		beginText := "{{BEGIN " + placeholder + "}}"
		endText := "{{END " + placeholder + "}}"

		origLines := strings.Split(string(fileContent), "\n")
		finalLines := []string{}
		for _, originalLine := range origLines {
			trimmedLine := strings.TrimSpace(originalLine)

			if strings.Contains(trimmedLine, endText) {
				isBeginPlaceholder = false
				finalLines = append(finalLines, globalSnippetIndentSpaces+trimmedLine)
				continue
			}

			if strings.Contains(trimmedLine, beginText) {
				isBeginPlaceholder = true
				finalLines = append(finalLines, originalLine)

				//Indent all the code with the same indentation as the 'begin' placeholder line
				prefixSpacesCount := len(originalLine) - len(strings.TrimLeft(originalLine, " "))
				globalSnippetIndentSpaces = strings.Repeat(" ", prefixSpacesCount)

				contentLines := strings.Split(content, "\n")
				for index, _ := range contentLines {
					contentLines[index] = globalSnippetIndentSpaces + contentLines[index]
				}
				finalLines = append(finalLines, strings.Join(contentLines, "\n"))

				logger.Debug("Appended content after line containing '%s'", beginText)

				continue
			}

			if isBeginPlaceholder {
				//Skip when between the placeholder tags
				continue
			}

			finalLines = append(finalLines, originalLine)
		}

		finalContent = strings.Join(finalLines, "\n")
	}
	err := ioutil.WriteFile(filePath, []byte(finalContent), 0600)
	CheckError(err)
}
