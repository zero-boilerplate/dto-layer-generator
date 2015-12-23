package plugins

import (
	"bytes"
	"fmt"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"strings"
	"text/template"

	"github.com/zero-boilerplate/dto-layer-generator/helpers"
	"github.com/zero-boilerplate/dto-layer-generator/setup"
)

func NewJavaPlugin() Plugin {
	p := &javaPlugin{}

	p.tpl = template.Must(template.New("name").Funcs(template.FuncMap{
		"print_field":       p.javaPrintFieldFunc,
		"print_field_class": p.javaPrintFieldClassFunc,
	}).Parse(`
		//region generated {{.Name}}
		// Generated with github.com/zero-boilerplate/dto-layer-generator
		private class {{.Name}} {
			{{range .Fields}}{{. | print_field}}
			{{end}}
			{{range .Fields}}{{. | print_field_class}}
			{{end}}
		}
		//endregion
	`))

	p.typeNameMap = map[string]string{
		"string":  "String",
		"bool":    "Boolean",
		"byte":    "Byte",
		"float32": "Float",
		"float64": "Float",
		"int":     "Integer",
		"int8":    "Integer",
		"int16":   "Integer",
		"int32":   "Integer",
		"int64":   "Integer",
		"uint":    "Integer",
		"uint8":   "Integer",
		"uint16":  "Integer",
		"uint32":  "Integer",
		"uint64":  "Integer",
	}

	return p
}

type javaPlugin struct {
	tpl         *template.Template
	typeNameMap map[string]string
}

func (j *javaPlugin) GenerateCode(logger helpers.Logger, dtoSetup *setup.DTOSetup) []byte {
	var outputBuf bytes.Buffer
	err := j.tpl.Execute(&outputBuf, dtoSetup)
	CheckError(err)

	return helpers.PrettifyCode(outputBuf.Bytes(), &helpers.PrettifyRules{
		MustPrefixWithEmptyLine:  func(trimmedLine string) bool { return strings.HasPrefix(trimmedLine, "private class") },
		StartIndentNextLine:      func(trimmedLine string) bool { return strings.HasSuffix(trimmedLine, "}") },
		StopIndentingCurrentLine: func(trimmedLine string) bool { return strings.HasSuffix(trimmedLine, "{") },
	})
}

func (j *javaPlugin) javaPrintFieldFunc(field *setup.DTOField) string {
	if field.IsObject() {
		return fmt.Sprintf("public %sClass %s;", field.Name, field.Name)
	}

	if field.IsObjectArray() {
		return fmt.Sprintf("public ArrayList<%sClass> %s;", field.Name, field.Name)
	}

	return fmt.Sprintf("public %s %s;", field.ConvertTypeName(j.typeNameMap), field.Name)
}

func (j *javaPlugin) javaPrintFieldClassFunc(field *setup.DTOField) string {
	if field.IsObject() || field.IsObjectArray() {
		childFields := []string{}
		for _, cf := range field.Fields {
			childFields = append(childFields, j.javaPrintFieldFunc(cf))
		}

		childFieldClasses := []string{}
		for _, cfc := range field.Fields {
			childFieldClasses = append(childFieldClasses, j.javaPrintFieldClassFunc(cfc))
		}

		return fmt.Sprintf(`private class %sClass {
			%s

			%s
		}`, field.Name, strings.Join(childFields, "\n"), strings.Join(childFieldClasses, "\n"))
	}
	return ""
}

func init() {
	RegisterPlugin("java", NewJavaPlugin())
}
