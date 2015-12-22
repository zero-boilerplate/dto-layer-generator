package setup

type DTOSetup struct {
	Name   string
	Url    string
	Fields []*DTOField
}

type DTOField struct {
	Name   string
	Type   string
	Fields []*DTOField
}
