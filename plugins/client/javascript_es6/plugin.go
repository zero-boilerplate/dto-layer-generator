package javascript_es6

import (
	"bytes"
	"fmt"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"net/url"
	"strings"
	"text/template"

	"github.com/zero-boilerplate/dto-layer-generator/helpers"
	"github.com/zero-boilerplate/dto-layer-generator/plugins"
	"github.com/zero-boilerplate/dto-layer-generator/setup"
)

func newPlugin() plugins.Plugin {
	p := &plugin{}

	p.tpl = template.Must(template.New("name").Funcs(template.FuncMap{
		"class_name_suffix":              p.funcClassNameSuffix,
		"join_field_names_for_url_query": p.funcJoinFieldNamesForUrlQuery,
		"constructor_params":             p.funcConstructorParams,
		"post_list":                      p.funcPostList,
		"patch_list":                     p.funcPatchList,
	}).Parse(`
		//region generated {{.Name}}
		// Generated with github.com/zero-boilerplate/dto-layer-generator

        {{$outerScope := .}}
        get{{$outerScope.Name}}URL() {
            return this.urlHelper.getApiV1Url('{{$outerScope.UrlWithStartingSlash}}');
        }

        {{if .IsInsertMethodEnabled}}
        insert{{$outerScope.Name}}({{.InsertableFields | constructor_params}}) {
            let url = this.get{{$outerScope.Name}}URL();
            let payload = {
                {{.InsertableFields | post_list}}
            };
            return this.http.post(url, payload);
        }
        {{end}}

        {{if .IsPatchMethodEnabled}}
        _newPatchReplaceMap(fieldName, value) {
            return {
                op: "replace",
                path: "/" + fieldName,
                value: value
            };
        }
        {{range .PatchableFieldGroups}}
        patch{{$outerScope.Name}}_{{. | class_name_suffix}}(id, {{. | constructor_params}}) {
            let url = this.get{{$outerScope.Name}}URL();
            url += "/" + id;

            let patchList = [
                {{. | patch_list}}
            ];

            return this.http.patch(url, patchList);
        }
        {{end}}
        {{end}}

        {{if .IsListMethodEnabled}}
        {{range .ListableFieldGroups}}
        list{{$outerScope.Name}}s_{{. | class_name_suffix}}() {
            let url = this.get{{$outerScope.Name}}URL();
            url += "?fields={{. | join_field_names_for_url_query}}";
            return this.http.get(url);
        }
        {{end}}
        {{end}}

        {{if .IsCountMethodEnabled}}
        countAll{{.Name}}s() {
            let url = this.urlHelper.getApiV1Url('{{.UrlWithStartingSlash}}');
            url += "?only_count_all=true";
            return this.http.get(url);
        }
        {{end}}

        {{if .IsGetMethodEnabled}}
        {{range .GetableFieldGroups}}
        get{{$outerScope.Name}}ById_{{. | class_name_suffix}}(id) {
            let url = this.get{{$outerScope.Name}}URL();
            url += "/" + id;
            url += "?fields={{. | join_field_names_for_url_query}}";
            return this.http.get(url);
        }
        {{end}}
        {{end}}

        {{if .IsDeleteMethodEnabled}}
        delete{{$outerScope.Name}}ById(id) {
            let url = this.get{{$outerScope.Name}}URL();
            url += "/" + id;
            return this.http.delete(url);
        }
        {{end}}

		//endregion
	`))

	return p
}

type plugin struct {
	tpl *template.Template
}

func (p *plugin) getUpfieldValueVariableString(field *setup.DTOField) string {
	fieldValueVar := field.NameToLowerCamelCase()

	if field.IsNumberUptype() {
		//Cast to Number in javascript
		return fmt.Sprintf("Number(%s)", fieldValueVar)
	} else {
		return fieldValueVar
	}
}

func (p *plugin) funcClassNameSuffix(dtoFields []*setup.DTOField) string {
	fieldNames := []string{}
	for _, f := range dtoFields {
		fieldNames = append(fieldNames, f.Name)
	}
	return strings.Join(fieldNames, "_")
}

func (p *plugin) funcJoinFieldNamesForUrlQuery(dtoFields []*setup.DTOField) string {
	lowercaseFieldNames := []string{}
	for _, f := range dtoFields {
		lowercaseFieldNames = append(lowercaseFieldNames, f.NameToSnakeCase())
	}
	encoded := url.QueryEscape(strings.Join(lowercaseFieldNames, "."))
	return encoded
}

func (p *plugin) funcConstructorParams(dtoFields []*setup.DTOField) string {
	param := []string{}
	for _, field := range dtoFields {
		param = append(param, field.NameToLowerCamelCase())
	}
	return strings.Join(param, ", ")
}

func (p *plugin) funcPostList(dtoFields []*setup.DTOField) string {
	line := []string{}
	for _, field := range dtoFields {
		line = append(line, fmt.Sprintf(`%s: "%s"`, field.Name, p.getUpfieldValueVariableString(field)))
	}
	return strings.Join(line, ",\n")
}

func (p *plugin) funcPatchList(dtoFields []*setup.DTOField) string {
	line := []string{}
	for _, field := range dtoFields {
		line = append(line, fmt.Sprintf(`this._newPatchReplaceMap("%s", %s)`, field.Name, p.getUpfieldValueVariableString(field)))
	}
	return strings.Join(line, ",\n")
}

func (p *plugin) GenerateCode(logger helpers.Logger, dtoSetup *setup.DTOSetup) []byte {
	var outputBuf bytes.Buffer
	err := p.tpl.Execute(&outputBuf, dtoSetup)
	CheckError(err)

	return helpers.PrettifyCode(outputBuf.Bytes(), &helpers.PrettifyRules{
		//MustPrefixWithEmptyLine:  func(trimmedLine string) bool { return strings.HasPrefix(trimmedLine, "private class") },
		StartIndentNextLine: func(trimmedLine string) bool {
			return strings.Count(trimmedLine, "{") > strings.Count(trimmedLine, "}") ||
				strings.Count(trimmedLine, "[") > strings.Count(trimmedLine, "]")
		},
		StopIndentingCurrentLine: func(trimmedLine string) bool {
			return strings.Count(trimmedLine, "{") < strings.Count(trimmedLine, "}") ||
				strings.Count(trimmedLine, "[") < strings.Count(trimmedLine, "]")
		},
	})
}

func init() {
	plugins.RegisterPlugin("client__javascript_es6", newPlugin())
}
