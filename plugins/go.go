package plugins

import (
	"bytes"
	"fmt"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"go/format"
	"strings"
	"text/template"

	"github.com/francoishill/dto-layer-generator/setup"
)

type GoPlugin struct{}

func (g *GoPlugin) GenerateCode(dtoSetup *setup.DTOSetup) []byte {
	var outputBuf bytes.Buffer
	err := golangTpl.Execute(&outputBuf, dtoSetup)
	CheckError(err)
	generatedGoCode := outputBuf.String()

	formattedGoCodeBytes, err := format.Source([]byte(generatedGoCode))
	CheckError(err)
	return formattedGoCodeBytes
}

func printFieldFunc(fieldIn interface{}) string {
	field := fieldIn.(*setup.DTOField)
	isObject := field.Type == "object"
	isObjectArray := field.Type == "objectarray"
	if isObject || isObjectArray {
		childFields := []string{}
		for _, cf := range field.Fields {
			childFields = append(childFields, printFieldFunc(cf))
		}

		structKeywordPrefix := ""
		if isObjectArray {
			structKeywordPrefix = "[]"
		}

		return fmt.Sprintf(`%s %sstruct {
			%s
		}`, field.Name, structKeywordPrefix, strings.Join(childFields, "\n"))
	}
	return fmt.Sprintf("%s %s", field.Name, field.Type)
}

var golangTpl = template.Must(template.New("name").Funcs(template.FuncMap{"print_field": printFieldFunc}).Parse(`
	package out

	type {{.Name}} struct {
		{{range .Fields}}{{. | print_field}}
		{{end}}
	}
`))
