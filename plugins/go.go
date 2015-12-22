package plugins

import (
	"bytes"
	"fmt"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"go/format"
	"strings"
	"text/template"

	"github.com/francoishill/dto-layer-generator/helpers"
	"github.com/francoishill/dto-layer-generator/setup"
)

func NewGoPlugin() Plugin {
	p := &goPlugin{}
	p.tpl = template.Must(template.New("name").Funcs(template.FuncMap{"print_field": p.goPrintFieldFunc}).Parse(`
		// Generated with github.com/francoishill/dto-layer-generator
		type {{.Name}} struct {
			{{range .Fields}}{{. | print_field}}
			{{end}}
		}
	`))
	return p
}

type goPlugin struct {
	tpl *template.Template
}

func (g *goPlugin) GenerateCode(logger Logger, dtoSetup *setup.DTOSetup) []byte {
	var outputBuf bytes.Buffer
	err := g.tpl.Execute(&outputBuf, dtoSetup)
	CheckError(err)

	prettyCode := helpers.PrettifyCode(outputBuf.Bytes(), &helpers.PrettifyRules{
		MustPrefixWithEmptyLine:  func(trimmedLine string) bool { return strings.HasSuffix(trimmedLine, "struct {") },
		StartIndentNextLine:      func(trimmedLine string) bool { return strings.HasSuffix(trimmedLine, "}") },
		StopIndentingCurrentLine: func(trimmedLine string) bool { return strings.HasSuffix(trimmedLine, "{") },
	})

	formattedCodeBytes, err := format.Source(prettyCode)
	if err != nil {
		logger.Warn("Unable to format (gofmt) the golang code, error was: %s", err.Error())
		return prettyCode
	}

	return formattedCodeBytes
}

func (g *goPlugin) goPrintFieldFunc(field *setup.DTOField) string {
	if field.IsObject() || field.IsObjectArray() {
		childFields := []string{}
		for _, cf := range field.Fields {
			childFields = append(childFields, g.goPrintFieldFunc(cf))
		}

		structKeywordPrefix := ""
		if field.IsObjectArray() {
			structKeywordPrefix = "[]"
		}

		return fmt.Sprintf(`%s %sstruct {
			%s
		}`, field.Name, structKeywordPrefix, strings.Join(childFields, "\n"))
	}
	return fmt.Sprintf("%s %s", field.Name, field.Type)
}

func init() {
	RegisterPlugin("go", NewGoPlugin())
}
