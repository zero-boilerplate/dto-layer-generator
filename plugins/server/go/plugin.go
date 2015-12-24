package server

import (
	"bytes"
	"fmt"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"go/format"
	"net/url"
	"strings"
	"text/template"

	"github.com/zero-boilerplate/dto-layer-generator/helpers"
	"github.com/zero-boilerplate/dto-layer-generator/plugins"
	"github.com/zero-boilerplate/dto-layer-generator/setup"
)

func newGoPlugin() plugins.Plugin {
	p := &goPlugin{}
	p.tpl = template.Must(template.New("name").Funcs(template.FuncMap{
		"fielddefs":                      p.funcFieldDefinitions,
		"class_name_suffix":              p.funcClassNameSuffix,
		"join_field_names_for_url_query": p.funcJoinFieldNamesForUrlQuery,
	}).Parse(`
		{{$outerScope := .}}
		// Generated with github.com/zero-boilerplate/dto-layer-generator
		
		func (c *controller) RelativeURLPatterns() []string {
			return []string{"{{.Url}}", "{{.Url}}/{id}"}
		}
		
		{{if .IsInsertMethodEnabled}}
		{{end}}
		
		{{if .IsPatchMethodEnabled}}
		{{end}}
		
		{{if .IsGetORListMethodEnabled}}

		
		{{if .IsGetMethodEnabled}}
		{{end}}
		{{if .IsListMethodEnabled}}
		{{range $group := .ListableFieldGroups}}
		type listEntity_{{. | class_name_suffix}} struct {
			{{$group | fielddefs}}
		}
		type listDTO_{{. | class_name_suffix}} struct {
			List       []*listEntity_{{. | class_name_suffix}}
			TotalCount int64
		}
		{{end}}
		{{end}}

		{{if .IsDeleteMethodEnabled}}
		func (c *controller) Delete(w http.ResponseWriter, r *http.Request) {
			authUser := c.GetUserFromRequest(r)
			id := c.MustUrlParamValue(r, "id").Int64()
			c.delete{{$outerScope.Name}}(authUser, id)
		}
		{{end}}

		func (c *controller) Get(w http.ResponseWriter, r *http.Request) {
			idVal := c.OptionalUrlParamValue(r, "id")
			
			authUser := c.GetUserFromRequest(r)

			if idVal.HasValue() {
				{{if .IsGetMethodEnabled}}
				//Get
				c.RenderJson(w, c.get{{$outerScope.Name}}ByIdRequestValue(authUser, idVal))
				{{else}}
				//Get method disabled in dto generator
				{{end}}
			} else {
				{{if .IsListMethodEnabled}}
				//List
				
				offset := c.OptionalQueryValue(r, "offset").Int64()

				limit := int64(DefaultNotificationLimit)
				if limitVal := c.OptionalQueryValue(r, "limit"); limitVal.HasValue() {
					limit = limitVal.Int64()
				}

				fields := c.MustQueryValue(r, "fields").String()
				switch fields {
				{{range $group := .ListableFieldGroups}}
				case "{{. | join_field_names_for_url_query}}":
					
					list := c.list{{$outerScope.Name}}s_{{. | class_name_suffix}}(authUser, offset, limit)
					totalCount := c.countAll{{$outerScope.Name}}s(authUser)

					entities := []*listEntity_{{. | class_name_suffix}}{}
					for _, e := range list {
						entities = append(entities, &listEntity_{{. | class_name_suffix}}{
							{{range $field := $group}}
							e.{{$field.Name}}(),
							{{end}}
						})
					}
					c.RenderJson(w, &listDTO_{{. | class_name_suffix}} {
						entities,
						totalCount,
					})

					break
				{{end}}
				default:
					panic(c.CreateHttpStatusClientError_BadRequest("Unsupported field combination: " + fields))
				}

				{{else}}
				//List method disabled in dto generator
				{{end}}
			}
		}
		{{end}}
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
		MustPrefixWithEmptyLine: func(trimmedLine string) bool { return strings.HasSuffix(trimmedLine, "struct {") },
		StartIndentNextLine: func(trimmedLine string) bool {
			return strings.Count(trimmedLine, "{") > strings.Count(trimmedLine, "}")
		},
		StopIndentingCurrentLine: func(trimmedLine string) bool {
			return strings.Count(trimmedLine, "{") < strings.Count(trimmedLine, "}")
		},
	})

	formattedCodeBytes, err := format.Source(prettyCode)
	if err != nil {
		logger.Warn("Unable to format (gofmt) the golang code, error was: %s", err.Error())
		return prettyCode
	}

	return formattedCodeBytes
}

func (g *goPlugin) funcFieldDefinitions(dtoFields []*setup.DTOField) string {
	lines := []string{}
	for _, field := range dtoFields {
		lines = append(lines, fmt.Sprintf(`%s %s`, field.Name, field.Type))
	}
	return strings.Join(lines, "\n")
}

func (g *goPlugin) funcClassNameSuffix(dtoFields []*setup.DTOField) string {
	fieldNames := []string{}
	for _, f := range dtoFields {
		fieldNames = append(fieldNames, f.Name)
	}
	return strings.Join(fieldNames, "_")
}

func (g *goPlugin) funcJoinFieldNamesForUrlQuery(dtoFields []*setup.DTOField) string {
	//TODO: Could make this a helper function as also used in plugin client/java_android?
	lowercaseFieldNames := []string{}
	for _, f := range dtoFields {
		lowercaseFieldNames = append(lowercaseFieldNames, f.NameToSnakeCase())
	}
	encoded := url.QueryEscape(strings.Join(lowercaseFieldNames, "."))
	return encoded
}

func init() {
	plugins.RegisterPlugin("go", newGoPlugin())
}
