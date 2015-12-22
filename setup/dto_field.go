package setup

import (
	"fmt"
)

type DTOField struct {
	Name   string
	Type   string
	Fields []*DTOField
}

func (d *DTOField) IsObject() bool      { return d.Type == "object" }
func (d *DTOField) IsObjectArray() bool { return d.Type == "objectarray" }

func (d *DTOField) ConvertTypeName(typeNameMap map[string]string) string {
	if converterTypeName, ok := typeNameMap[d.Type]; ok {
		return converterTypeName
	}
	panic(fmt.Sprintf("Cannot convert TypeName from %s using map: %#v", d.Type, typeNameMap))
}
