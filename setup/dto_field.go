package setup

import (
	"fmt"

	"github.com/zero-boilerplate/dto-layer-generator/helpers"
)

type DTOField struct {
	Name string
	Type string
}

func (d *DTOField) ConvertTypeName(typeNameMap map[string]string) string {
	if converterTypeName, ok := typeNameMap[d.Type]; ok {
		return converterTypeName
	}
	panic(fmt.Sprintf("Cannot convert TypeName from %s using map: %#v", d.Type, typeNameMap))
}

func (d *DTOField) NameToLowerCamelCase() string {
	return helpers.ToLowerCamelCase(d.Name)
}

func (d *DTOField) NameToKebabCase() string {
	return helpers.ToKebabCase(d.Name)
}

func (d *DTOField) NameToSnakeCase() string {
	return helpers.ToSnakeCase(d.Name)
}
