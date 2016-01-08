package java_retrofit

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
		"up_field_type_name":             p.funcUpFieldTypeName,
		"down_field_type_name":           p.funcDownFieldTypeName,
		"up_fielddefs":                   p.funcUpfieldDefinitions,
		"down_fielddefs":                 p.funcDownfieldDefinitions,
		"up_constructor_params":          p.funcUpConstructorParams,
		"request_constructor_body":       p.funcRequestConstructorBody,
		"up_single_fielddef":             p.funcSingleUpFieldDefinition,
		"class_name_suffix":              p.funcClassNameSuffix,
		"patch_request_constructor_body": p.funcPatchRequestConstructorBody,
		"join_field_names_for_url_query": p.funcJoinFieldNamesForUrlQuery,
	}).Parse(`
		//region generated {{.Name}}
		// Generated with github.com/zero-boilerplate/dto-layer-generator

        {{$outerScope := .}}
        {{if .IsInsertMethodEnabled}}
		private static class {{.Name}}InsertDTO {
			private static class Request {
	            {{.InsertableFields | up_fielddefs}}

	            private Request({{.InsertableFields | up_constructor_params}}) {
	                {{.InsertableFields | request_constructor_body}}
	            }
	        }

	        private static class Response {
	            {{.IdField | up_single_fielddef}}
	        }			
		}
        {{end}}

        {{if .IsPatchMethodEnabled}}
		private static class {{.Name}}PatchDTOs {
			private static HashMap<String, Object> newReplaceMap(String fieldName, Object value) {
	            //JSON PATCH Format: rfc6902 -- http://tools.ietf.org/html/rfc6902 -- http://jsonpatch.com/
	            HashMap<String, Object> map = new HashMap<>();
	            map.put("op", "replace");
	            map.put("path", "/" + fieldName);
	            map.put("value", value);
	            return map;
	        }

	        {{range .PatchableFieldGroups}}
	        private static class Request_{{. | class_name_suffix}} extends ArrayList<HashMap<String, Object>> {
	            private Request_{{. | class_name_suffix}}({{. | up_constructor_params}}) {
	                super();
	                {{. | patch_request_constructor_body}}
	            }
	        }
            {{end}}
        }
        {{end}}

        {{if .IsListMethodEnabled}}
        private static class {{.Name}}ListDTOs {
            {{range .ListableFieldGroups}}
            private static class Response_{{. | class_name_suffix}} {
                private static class ListItem {
                    {{. | down_fielddefs}}
                }

                public ArrayList<ListItem> List;
                public Integer TotalCount;
            }
            {{end}}
        }
        {{end}}

        {{if .IsGetMethodEnabled}}
        private static class {{.Name}}GetDTOs {
            {{range .GetableFieldGroups}}
            private static class Response_{{. | class_name_suffix}} {
                {{. | down_fielddefs}}
            }
            {{end}}
        }
        {{end}}

        private interface I{{.Name}}WebService {
            //This requires Retrofit to be references: http://square.github.io/retrofit/

            {{if .IsInsertMethodEnabled}}
            // Insert
            @POST("{{$outerScope.Url}}")
            Call<{{$outerScope.Name}}InsertDTO.Response> insert(@Body {{$outerScope.Name}}InsertDTO.Request body);
            {{end}}

            {{if .IsPatchMethodEnabled}}
            // Patch/update
            {{range .PatchableFieldGroups}}
            @PATCH("{{$outerScope.Url}}/{id}")
            Call<Void> patch(@Path("id") {{$outerScope.IdField | up_field_type_name}} id, @Body {{$outerScope.Name}}PatchDTOs.Request_{{. | class_name_suffix}} body);
            {{end}}
            {{end}}

            {{if .IsListMethodEnabled}}
            // List
            {{range .ListableFieldGroups}}
            @GET("{{$outerScope.Url}}?fields={{. | join_field_names_for_url_query}}")
            Call<{{$outerScope.Name}}ListDTOs.Response_{{. | class_name_suffix}}> list_{{. | class_name_suffix}}();
            {{end}}
            {{end}}
            
            {{if .IsCountMethodEnabled}}
            //Count
            @GET("{{$outerScope.Url}}?only_count_all=true")
            Call<Integer> countAll{{.Name}}s();
            {{end}}

            {{if .IsGetMethodEnabled}}
            // Get single
            {{range .GetableFieldGroups}}
            @GET("{{$outerScope.Url}}/{id}?fields={{. | join_field_names_for_url_query}}")
            Call<{{$outerScope.Name}}GetDTOs.Response_{{. | class_name_suffix}}> get_{{. | class_name_suffix}}(@Path("id") {{$outerScope.IdField | down_field_type_name}} id);
            {{end}}
            {{end}}

            {{if .IsDeleteMethodEnabled}}
            @DELETE("{{$outerScope.Url}}/{id}")
            Call<Void> delete(@Path("id") {{$outerScope.IdField | up_field_type_name}} id);
            {{end}}
        }

		//endregion
	`))

	p.typeNameMap = map[string]string{
		"string":    "String",
		"bool":      "Boolean",
		"byte":      "Byte",
		"float32":   "Float",
		"float64":   "Float",
		"int":       "Integer",
		"int8":      "Integer",
		"int16":     "Integer",
		"int32":     "Integer",
		"int64":     "Integer",
		"uint":      "Integer",
		"uint8":     "Integer",
		"uint16":    "Integer",
		"uint32":    "Integer",
		"uint64":    "Integer",
		"time.Time": "Date",
	}

	return p
}

type plugin struct {
	tpl         *template.Template
	typeNameMap map[string]string
}

func (p *plugin) funcUpFieldTypeName(dtoField *setup.DTOField) string {
	return dtoField.ConvertUptypeName(p.typeNameMap)
}

func (p *plugin) funcDownFieldTypeName(dtoField *setup.DTOField) string {
	return dtoField.ConvertDowntypeName(p.typeNameMap)
}

func (p *plugin) funcUpfieldDefinitions(dtoFields []*setup.DTOField) string {
	lines := []string{}
	for _, field := range dtoFields {
		lines = append(lines, fmt.Sprintf(`public %s %s;`, field.ConvertUptypeName(p.typeNameMap), field.Name))
	}
	return strings.Join(lines, "\n")
}

func (p *plugin) funcDownfieldDefinitions(dtoFields []*setup.DTOField) string {
	lines := []string{}
	for _, field := range dtoFields {
		lines = append(lines, fmt.Sprintf(`public %s %s;`, field.ConvertDowntypeName(p.typeNameMap), field.Name))
	}
	return strings.Join(lines, "\n")
}

func (p *plugin) funcUpConstructorParams(dtoFields []*setup.DTOField) string {
	param := []string{}
	for _, field := range dtoFields {
		param = append(param, fmt.Sprintf(`%s %s`, field.ConvertUptypeName(p.typeNameMap), field.NameToLowerCamelCase()))
	}
	return strings.Join(param, ", ")
}

func (p *plugin) funcRequestConstructorBody(dtoFields []*setup.DTOField) string {
	line := []string{}
	for _, field := range dtoFields {
		line = append(line, fmt.Sprintf(`this.%s = %s;`, field.Name, field.NameToLowerCamelCase()))
	}
	return strings.Join(line, "\n")
}

func (p *plugin) funcSingleUpFieldDefinition(dtoField *setup.DTOField) string {
	return p.funcUpfieldDefinitions([]*setup.DTOField{dtoField})
}

func (p *plugin) funcClassNameSuffix(dtoFields []*setup.DTOField) string {
	fieldNames := []string{}
	for _, f := range dtoFields {
		fieldNames = append(fieldNames, f.Name)
	}
	return strings.Join(fieldNames, "_")
}

func (p *plugin) funcPatchRequestConstructorBody(dtoFields []*setup.DTOField) string {
	line := []string{}
	for _, field := range dtoFields {
		line = append(line, fmt.Sprintf(`add(newReplaceMap("%s", %s));`, field.Name, field.NameToLowerCamelCase()))
	}
	return strings.Join(line, "\n")
}

func (p *plugin) funcJoinFieldNamesForUrlQuery(dtoFields []*setup.DTOField) string {
	lowercaseFieldNames := []string{}
	for _, f := range dtoFields {
		lowercaseFieldNames = append(lowercaseFieldNames, f.NameToSnakeCase())
	}
	encoded := url.QueryEscape(strings.Join(lowercaseFieldNames, "."))
	return encoded
}

func (p *plugin) GenerateCode(logger helpers.Logger, dtoSetup *setup.DTOSetup) []byte {
	var outputBuf bytes.Buffer
	err := p.tpl.Execute(&outputBuf, dtoSetup)
	CheckError(err)

	return helpers.PrettifyCode(outputBuf.Bytes(), &helpers.PrettifyRules{
		MustPrefixWithEmptyLine: func(trimmedLine string) bool { return strings.HasPrefix(trimmedLine, "private class") },
		StartIndentNextLine: func(trimmedLine string) bool {
			return strings.Count(trimmedLine, "{") > strings.Count(trimmedLine, "}")
		},
		StopIndentingCurrentLine: func(trimmedLine string) bool {
			return strings.Count(trimmedLine, "{") < strings.Count(trimmedLine, "}")
		},
	})
}

func init() {
	plugins.RegisterPlugin("client__java_retrofit", newPlugin())
}
