package helpers

import (
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"github.com/francoishill/golang-web-dry/osutils"
	"io/ioutil"
	"strings"
)

func InjectContentIntoFilePlaceholder(filePath, placeholder, content string) {
	finalContent := content
	if osutils.FileExists(filePath) {
		fileContent, err := ioutil.ReadFile(filePath)
		CheckError(err)

		isBeginPlaceholder := false

		beginText := "{{BEGIN " + placeholder + "}}"
		endText := "{{END " + placeholder + "}}"

		origLines := strings.Split(string(fileContent), "\n")
		finalLines := []string{}
		for _, line := range origLines {
			trimmedLine := strings.TrimSpace(line)

			if strings.Contains(trimmedLine, endText) {
				isBeginPlaceholder = false
				finalLines = append(finalLines, trimmedLine)
				continue
			}

			if strings.Contains(trimmedLine, beginText) {
				isBeginPlaceholder = true
				finalLines = append(finalLines, trimmedLine)
				finalLines = append(finalLines, content)
				continue
			}

			if isBeginPlaceholder {
				//Skip when between the placeholder tags
				continue
			}

			finalLines = append(finalLines, trimmedLine)
		}
		finalContent = strings.Join(finalLines, "\n")
	}
	err := ioutil.WriteFile(filePath, []byte(finalContent), 0600)
	CheckError(err)
}
