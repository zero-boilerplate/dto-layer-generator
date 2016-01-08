package setup

import (
	"fmt"

	"github.com/zero-boilerplate/dto-layer-generator/helpers"
)

type DTOField struct {
	Name     string
	Type     string
	Uptype   string
	Downtype string
}

func (d *DTOField) ConvertUptypeName(typeNameMap map[string]string) string {
	if converterTypeName, ok := typeNameMap[d.Uptype]; ok {
		return converterTypeName
	}
	panic(fmt.Sprintf("Cannot convert TypeName from %s using map: %#v", d.Uptype, typeNameMap))
}

func (d *DTOField) ConvertDowntypeName(typeNameMap map[string]string) string {
	if converterTypeName, ok := typeNameMap[d.Downtype]; ok {
		return converterTypeName
	}
	panic(fmt.Sprintf("Cannot convert TypeName from %s using map: %#v", d.Downtype, typeNameMap))
}

func fieldIsNumberType(fieldType string) bool {
	switch fieldType {
	case "float64", "float32", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "byte":
		return true
	default:
		return false
	}
}

func (d *DTOField) IsNumberUptype() bool {
	return fieldIsNumberType(d.Uptype)
}

func (d *DTOField) IsNumberDowntype() bool {
	return fieldIsNumberType(d.Downtype)
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
