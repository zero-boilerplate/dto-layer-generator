package server

import (
	"bytes"
	"fmt"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"go/format"
	"strings"
	"text/template"

	"github.com/zero-boilerplate/dto-layer-generator/helpers"
	"github.com/zero-boilerplate/dto-layer-generator/plugins"
	"github.com/zero-boilerplate/dto-layer-generator/setup"
)

func newGoPlugin() plugins.Plugin {
	p := &goPlugin{}
	p.tpl = template.Must(template.New("name").Funcs(template.FuncMap{"print_field": p.goPrintFieldFunc}).Parse(`
		// Generated with github.com/zero-boilerplate/dto-layer-generator
		type {{.Name}} struct {
			{{range .AllFields}}{{. | print_field}}
			{{end}}
		}
	`))
	return p
}

type goPlugin struct {
	tpl *template.Template
}

func (g *goPlugin) GenerateCode(logger helpers.Logger, dtoSetup *setup.DTOSetup) []byte {
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
	return fmt.Sprintf("%s %s", field.Name, field.Type)
}

func init() {
	plugins.RegisterPlugin("go", newGoPlugin())
}
