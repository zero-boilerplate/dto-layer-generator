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

func (d *DTOField) IsNumberType() bool {
	switch d.Type {
	case "float64", "float32", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "byte":
		return true
	default:
		return false
	}
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
