package go_negroni

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

func newPlugin() plugins.Plugin {
	p := &plugin{}
	p.tpl = template.Must(template.New("name").Funcs(template.FuncMap{
		"field_type_name":                p.funcFieldTypeName,
		"fielddefs":                      p.funcFieldDefinitions,
		"class_name_suffix":              p.funcClassNameSuffix,
		"join_field_names_for_url_query": p.funcJoinFieldNamesForUrlQuery,
	}).Parse(`
		{{$outerScope := .}}
		// Generated with github.com/zero-boilerplate/dto-layer-generator
		// Required imports:
		// - "github.com/mholt/binding"
		
		func (c *controller) RelativeURLPatterns() []string {
			return []string{"{{.UrlWithStartingSlash}}", "{{.UrlWithStartingSlash}}/{id}"}
		}
		
		{{if .IsInsertMethodEnabled}}
		type insertDTO_Request struct {
			{{.InsertableFields | fielddefs}}
		}

		func (i *insertDTO_Request) FieldMap(r *http.Request) binding.FieldMap {
			return binding.FieldMap{}
		}
		{{end}}
		
		{{if .IsPatchMethodEnabled}}
		//JSON PATCH Format: rfc6902 -- http://tools.ietf.org/html/rfc6902 -- http://jsonpatch.com/
		type patchOperationDTO struct {
			Operation string      ` + "`json:\"op\"`" + `
			Path      string      ` + "`json:\"path\"`" + `
			Value     interface{} ` + "`json:\"value\"`" + `
		}
		{{end}}
		
		{{if .IsGetMethodEnabled}}
		{{range $group := .GetableFieldGroups}}
		type getDTO_Response__{{. | class_name_suffix}} struct {
			{{$group | fielddefs}}
		}
		{{end}}
		{{end}}

		{{if .IsListMethodEnabled}}
		{{range $group := .ListableFieldGroups}}
		type listEntity_Response__{{. | class_name_suffix}} struct {
			{{$group | fielddefs}}
		}
		type listDTO_{{. | class_name_suffix}} struct {
			List       []*listEntity_Response__{{. | class_name_suffix}}
			TotalCount int64
		}
		{{end}}
		{{end}}

		

		{{if .IsInsertMethodEnabled}}
		func (c *controller) Insert(w http.ResponseWriter, r *http.Request) {
			dto := &insertDTO_Request{}
			errs := binding.Json(r, dto)
			if errs.Handle(w) {
				return
			}

			authUser := c.GetUserFromRequest(r)
			id := c.insert{{$outerScope.Name}}(authUser, dto)
			c.RenderJson(w, struct{
				Id {{$outerScope.IdField | field_type_name}}
			}{id})
		}
		{{end}}

		{{if .IsPatchMethodEnabled}}
		func (c *controller) Patch(w http.ResponseWriter, r *http.Request) {
			idVal := c.MustUrlParamValue(r, "id")

			dbTx := c.MustBeginTransaction()
			defer c.DeferableCommitOnSuccessRollbackOnFail(dbTx)

			operations := []*patchOperationDTO{}
			c.Ctx.Misc.HttpRequestHelperService.DecodeJsonRequest(r, &operations)
			for _, o := range operations {
				switch o.Operation {
				case "replace":
					switch o.Path {
					{{range $field := .AllUniquePatchableFields}}
					case "/{{$field.Name}}":
						authUser := c.GetUserFromRequest(r)
						c.set{{$outerScope.Name}}Field_{{$field.Name}}_ById(dbTx, authUser, idVal, o.Value.({{$field.Type}}))
						break
					{{end}}
					default:
						panic(c.CreateHttpStatusClientError_BadRequest("Unsupported replace field name for '{{$outerScope.Name}}' entity: " + o.Path))
					}
					break
				default:
					panic(c.CreateHttpStatusClientError_BadRequest("Unsupported path operation type: " + o.Operation))
				}
			}
		}
		{{end}}

		{{if .IsDeleteMethodEnabled}}
		func (c *controller) Delete(w http.ResponseWriter, r *http.Request) {
			authUser := c.GetUserFromRequest(r)
			idVal := c.MustUrlParamValue(r, "id")
			c.delete{{$outerScope.Name}}ById(authUser, idVal)
		}
		{{end}}

		{{if .IsGetORListORCountMethodEnabled}}
		func (c *controller) Get(w http.ResponseWriter, r *http.Request) {
			idVal := c.OptionalUrlParamValue(r, "id")
			
			authUser := c.GetUserFromRequest(r)

			if idVal.HasValue() {
				{{if .IsGetMethodEnabled}}
				//Get

				fields := c.MustQueryValue(r, "fields").String()
				switch fields {
				{{range $group := .GetableFieldGroups}}
				case "{{. | join_field_names_for_url_query}}":					
					entity := c.get{{$outerScope.Name}}ById(authUser, idVal)
					c.RenderJson(w, &getDTO_Response__{{. | class_name_suffix}}{
						{{range $field := $group}}
						entity.{{$field.Name}}(),
						{{end}}
					})
					break
				{{end}}
				default:
					panic(c.CreateHttpStatusClientError_BadRequest("Unsupported field combination: " + fields))
				}

				{{else}}
				panic(c.CreateHttpStatusClientError_BadRequest("GET by ID disabled for {{$outerScope.Name}}"))
				{{end}}

				return
			}
				
			{{if .IsCountMethodEnabled}}
			//Count
			possibleOnlyCountAllValue := c.OptionalQueryValue(r, "only_count_all")
			if possibleOnlyCountAllValue.HasValue() && possibleOnlyCountAllValue.Bool() {
				totalCount := c.countAll{{$outerScope.Name}}s(authUser)
				c.RenderJson(w, totalCount)
				return
			}
			{{end}}

			{{if .IsListMethodEnabled}}
			//List
			offset := c.OptionalQueryValue(r, "offset").Int64()

			limit := int64(Default{{$outerScope.Name}}Limit)
			if limitVal := c.OptionalQueryValue(r, "limit"); limitVal.HasValue() {
				limit = limitVal.Int64()
			}

			fields := c.MustQueryValue(r, "fields").String()
			switch fields {
			{{range $group := .ListableFieldGroups}}
			case "{{. | join_field_names_for_url_query}}":
				
				list := c.list{{$outerScope.Name}}s_{{. | class_name_suffix}}(authUser, offset, limit)
				totalCount := c.countAll{{$outerScope.Name}}s(authUser)

				entities := []*listEntity_Response__{{. | class_name_suffix}}{}
				for _, e := range list {
					entities = append(entities, &listEntity_Response__{{. | class_name_suffix}}{
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
			panic(c.CreateHttpStatusClientError_BadRequest("List disabled for {{$outerScope.Name}}"))
			{{end}}			
		}
		{{end}}
	`))
	return p
}

type plugin struct {
	tpl *template.Template
}

func (p *plugin) GenerateCode(logger helpers.Logger, dtoSetup *setup.DTOSetup) []byte {
	var outputBuf bytes.Buffer
	err := p.tpl.Execute(&outputBuf, dtoSetup)
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

func (p *plugin) funcFieldTypeName(dtoField *setup.DTOField) string {
	return dtoField.Type
}

func (p *plugin) funcFieldDefinitions(dtoFields []*setup.DTOField) string {
	lines := []string{}
	for _, field := range dtoFields {
		lines = append(lines, fmt.Sprintf(`%s %s`, field.Name, field.Type))
	}
	return strings.Join(lines, "\n")
}

func (p *plugin) funcClassNameSuffix(dtoFields []*setup.DTOField) string {
	fieldNames := []string{}
	for _, f := range dtoFields {
		fieldNames = append(fieldNames, f.Name)
	}
	return strings.Join(fieldNames, "_")
}

func (p *plugin) funcJoinFieldNamesForUrlQuery(dtoFields []*setup.DTOField) string {
	//TODO: Could make this a helper function as also used in plugin client/java_android?
	lowercaseFieldNames := []string{}
	for _, f := range dtoFields {
		lowercaseFieldNames = append(lowercaseFieldNames, f.NameToSnakeCase())
	}
	encoded := url.QueryEscape(strings.Join(lowercaseFieldNames, "."))
	return encoded
}

func init() {
	plugins.RegisterPlugin("server__go_common_ddd", newPlugin())
}
